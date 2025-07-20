package main

import (
	"database/sql"
	"fmt"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/config"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/delivery/grpc_server"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/gateway/messaging"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/repository"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/usecase"
	"github.com/apriliantocecep/posfin-blog/services/auth/pkg/pb"
	"github.com/apriliantocecep/posfin-blog/shared"
	sharedlib "github.com/apriliantocecep/posfin-blog/shared/lib"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// vault client
	vaultClient := shared.NewVaultClient()
	secret := utils.GetVaultSecretConfig(vaultClient)

	// rabbitmq client
	rabbitMQClient := sharedlib.NewRabbitMQClient(vaultClient)
	defer func(Conn *amqp.Connection) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("error closing rabbitmq: %v", err)
		}
	}(rabbitMQClient.Conn)

	// setup publisher
	userCreatedPublisher := messaging.NewUserPublisher(rabbitMQClient.Conn, "user_created", "user_created")

	// dependencies
	database := config.NewDatabase(vaultClient)
	defer func(Conn *sql.DB) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("error closing db: %v", err)
		}
	}(database.Conn)
	jwt := config.NewJwt(vaultClient)
	userRepository := repository.NewUserRepository()
	userUseCase := usecase.NewUserUseCase(userRepository, jwt, database.DB, userCreatedPublisher)

	// grpc server
	srv := grpc_server.NewAuthServer(userUseCase)
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, srv)

	// listener
	portStr := secret["AUTH_SERVICE_PORT"]
	if portStr == nil || portStr == "" {
		log.Fatalln("AUTH_SERVICE_PORT is not set")
	}
	port := utils.ParsePort(portStr.(string))

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", listen.Addr())

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
