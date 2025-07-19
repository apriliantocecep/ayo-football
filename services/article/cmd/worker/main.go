package main

import (
	"context"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/config"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/delivery/messaging"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/repository"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/usecase"
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

	// dependencies
	database := config.NewDatabase(vaultClient)
	defer func() {
		ctx := context.Background()
		if err := database.Client.Disconnect(ctx); err != nil {
			log.Fatalf("error closing db: %v", err)
		}
	}()
	articleDb := database.Client.Database("posfin")
	articleCollection := articleDb.Collection("articles")
	articleRepository := repository.NewArticleRepository(articleCollection)
	moderationUseCase := usecase.NewModerationUseCase(database.Client, articleRepository)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// run consumer
	wg.Add(1) // wg.Add(total_consumer)
	go func() {
		defer wg.Done()
		log.Println("[worker] MetadataConsumer started")
		runMetadataConsumer(ctx, rabbitMQClient.Channel, moderationUseCase)
	}()

	// signal shutdown
	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	s := <-terminateSignals
	log.Printf("Shutdown signal received: %s", s)
	cancel()
	wg.Wait()
	log.Println("All consumers stopped gracefully.")
}

func runMetadataConsumer(ctx context.Context, channel *amqp.Channel, moderationUseCase *usecase.ModerationUseCase) {
	consumer := messaging.NewMetadataConsumer(moderationUseCase)
	sharedmessaging.ConsumeQueue(ctx, channel, "moderation_checker", consumer.Consume)
}
