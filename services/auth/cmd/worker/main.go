package main

import (
	"context"
	"database/sql"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/config"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/delivery/messaging"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/repository"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/usecase"
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
	wg.Add(1) // wg.Add(total_consumer)
	go func() {
		defer wg.Done()
		log.Println("[worker] UserCreatedConsumer started")
		runUserCreatedConsumer(ctx, rabbitMQClient.Conn, registerUseCase)
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
