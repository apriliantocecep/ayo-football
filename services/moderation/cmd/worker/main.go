package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/config"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/delivery/messaging"
	gatewaymessaging "github.com/apriliantocecep/posfin-blog/services/moderation/internal/gateway/messaging"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/repository"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/usecase"
	"github.com/apriliantocecep/posfin-blog/shared"
	sharedlib "github.com/apriliantocecep/posfin-blog/shared/lib"
	sharedmessaging "github.com/apriliantocecep/posfin-blog/shared/messaging"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// vault client
	vaultClient := shared.NewVaultClient()

	// consul client
	consul := sharedlib.NewConsulClient(vaultClient)

	// register service to consul
	ttl := time.Second * 8
	address := os.Getenv("SERVICE_URL")
	if address == "" {
		address = "127.0.0.1"
	}
	workerName := utils.Getenv("WORKER_NAME", uuid.NewString())
	serviceName := "moderation-worker-cluster"
	serviceRegisteredID := fmt.Sprintf("moderation-worker-%s", workerName)
	checkID := fmt.Sprintf("worker:%s", serviceRegisteredID)
	tags := []string{
		"traefik.enable=false",
		"role=worker",
	}
	consul.ServiceRegister(serviceRegisteredID, serviceName, address, 3000, &capi.AgentServiceCheck{
		TTL:                            ttl.String(),
		DeregisterCriticalServiceAfter: "1m",
		TLSSkipVerify:                  true,
		CheckID:                        checkID,
	}, tags)

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
	wg.Add(3) // wg.Add(total goroutine)
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
	go func() {
		defer wg.Done()
		log.Println("[worker] health check started")
		consul.StartTTLHeartbeat(checkID, time.Second*5)
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

func runArticleCreatedConsumer(ctx context.Context, rabbitMQConn *amqp.Connection, metadataUseCase *usecase.MetadataUseCase, moderationUseCase *usecase.ModerationUseCase) {
	consumer := messaging.NewArticleConsumer(metadataUseCase, moderationUseCase)
	sharedmessaging.ConsumeQueue(ctx, rabbitMQConn, "article_created", consumer.ConsumeArticleCreated)
}

func runArticleModerationConsumer(ctx context.Context, rabbitMQConn *amqp.Connection, metadataUseCase *usecase.MetadataUseCase, moderationUseCase *usecase.ModerationUseCase) {
	consumer := messaging.NewArticleConsumer(metadataUseCase, moderationUseCase)
	sharedmessaging.ConsumeQueue(ctx, rabbitMQConn, "article_moderation", consumer.ConsumeArticleModeration)
}
