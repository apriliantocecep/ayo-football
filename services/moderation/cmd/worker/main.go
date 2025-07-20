package main

import (
	"context"
	"database/sql"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/config"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/delivery/messaging"
	gatewaymessaging "github.com/apriliantocecep/posfin-blog/services/moderation/internal/gateway/messaging"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/repository"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/usecase"
	"github.com/apriliantocecep/posfin-blog/shared"
	sharedlib "github.com/apriliantocecep/posfin-blog/shared/lib"
	sharedmessaging "github.com/apriliantocecep/posfin-blog/shared/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// vault client
	vaultClient := shared.NewVaultClient()

	// rabbitmq client
	rabbitMQClient := sharedlib.NewRabbitMQClient(vaultClient)
	defer func(Conn *amqp.Connection) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("error closing rabbitmq: %v", err)
		}
	}(rabbitMQClient.Conn)

	// setup publisher
	metadataPublisher := gatewaymessaging.NewMetadataPublisher(rabbitMQClient.Conn, "moderation_checker", "moderation_checker")

	// dependencies
	database := config.NewDatabase(vaultClient)
	defer func(Conn *sql.DB) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("error closing db: %v", err)
		}
	}(database.Conn)
	metadataRepository := repository.NewMetadataRepository()
	metadataUseCase := usecase.NewMetadataUseCase(database.DB, metadataRepository)
	moderationUseCase := usecase.NewModerationUseCase(database.DB, metadataRepository, metadataPublisher)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// run consumer
	wg.Add(2) // wg.Add(total_consumer)
	go func() {
		defer wg.Done()
		log.Println("[worker] ArticleCreatedConsumer started")
		runArticleCreatedConsumer(ctx, rabbitMQClient.Conn, metadataUseCase, moderationUseCase)
	}()
	go func() {
		defer wg.Done()
		log.Println("[worker] ArticleModerationConsumer started")
		runArticleModerationConsumer(ctx, rabbitMQClient.Conn, metadataUseCase, moderationUseCase)
	}()

	// signal shutdown
	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	s := <-terminateSignals
	log.Printf("Shutdown signal received: %s", s)
	cancel()
	wg.Wait()
	log.Println("All consumers stopped gracefully.")

	//stop := false
	//for !stop {
	//	select {
	//	case s := <-terminateSignals:
	//		log.Printf("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME: %s", s)
	//		cancel()
	//		stop = true
	//	}
	//}
	//
	//time.Sleep(5 * time.Second) // wait for all consumers to finish processing
}

func runArticleCreatedConsumer(ctx context.Context, rabbitMQConn *amqp.Connection, metadataUseCase *usecase.MetadataUseCase, moderationUseCase *usecase.ModerationUseCase) {
	consumer := messaging.NewArticleConsumer(metadataUseCase, moderationUseCase)
	sharedmessaging.ConsumeQueue(ctx, rabbitMQConn, "article_created", consumer.ConsumeArticleCreated)
}

func runArticleModerationConsumer(ctx context.Context, rabbitMQConn *amqp.Connection, metadataUseCase *usecase.MetadataUseCase, moderationUseCase *usecase.ModerationUseCase) {
	consumer := messaging.NewArticleConsumer(metadataUseCase, moderationUseCase)
	sharedmessaging.ConsumeQueue(ctx, rabbitMQConn, "article_moderation", consumer.ConsumeArticleModeration)
}
