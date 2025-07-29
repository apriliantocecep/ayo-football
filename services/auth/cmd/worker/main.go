package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/config"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/delivery/messaging"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/repository"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/usecase"
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
	serviceName := "auth-worker-cluster"
	serviceRegisteredID := fmt.Sprintf("auth-worker-%s", workerName)
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
	defer func(Conn *sql.DB) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("error closing db: %v", err)
		}
	}(database.Conn)
	userRepository := repository.NewUserRepository()
	registerUseCase := usecase.NewRegisterUseCase(database.DB, userRepository)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// run consumer
	wg.Add(2) // wg.Add(total goroutine)
	go func() {
		defer wg.Done()
		log.Println("[worker] UserCreatedConsumer started")
		runUserCreatedConsumer(ctx, rabbitMQClient.Conn, registerUseCase)
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

func runUserCreatedConsumer(ctx context.Context, rabbitMQConn *amqp.Connection, registerUseCase *usecase.RegisterUseCase) {
	consumer := messaging.NewUserConsumer(registerUseCase)
	sharedmessaging.ConsumeQueue(ctx, rabbitMQConn, "user_created", consumer.ConsumeUserCreated)
}
