package main

import (
	"context"
	"fmt"
	"github.com/apriliantocecep/ayo-football/services/article/internal/config"
	"github.com/apriliantocecep/ayo-football/services/article/internal/delivery/messaging"
	"github.com/apriliantocecep/ayo-football/services/article/internal/repository"
	"github.com/apriliantocecep/ayo-football/services/article/internal/usecase"
	"github.com/apriliantocecep/ayo-football/shared"
	sharedlib "github.com/apriliantocecep/ayo-football/shared/lib"
	sharedmessaging "github.com/apriliantocecep/ayo-football/shared/messaging"
	"github.com/apriliantocecep/ayo-football/shared/utils"
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
	serviceName := "article-worker-cluster"
	serviceRegisteredID := fmt.Sprintf("article-worker-%s", workerName)
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
	wg.Add(2) // wg.Add(total goroutine)
	go func() {
		defer wg.Done()
		log.Println("[worker] MetadataConsumer started")
		runMetadataConsumer(ctx, rabbitMQClient.Conn, moderationUseCase)
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

func runMetadataConsumer(ctx context.Context, rabbitMQConn *amqp.Connection, moderationUseCase *usecase.ModerationUseCase) {
	consumer := messaging.NewMetadataConsumer(moderationUseCase)
	sharedmessaging.ConsumeQueue(ctx, rabbitMQConn, "moderation_checker", consumer.Consume)
}
