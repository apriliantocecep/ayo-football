package main

import (
	"context"
	"fmt"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/config"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/delivery/grpc_server"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/gateway/messaging"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/repository"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/usecase"
	"github.com/apriliantocecep/posfin-blog/services/article/pkg/pb"
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
	articleCreatedPublisher := messaging.NewArticlePublisher(rabbitMQClient.Channel, "article_created")

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
	articleUseCase := usecase.NewArticleUseCase(database.Client, articleRepository, articleCreatedPublisher)

	// grpc server
	srv := grpc_server.NewArticleServer(articleUseCase)
	s := grpc.NewServer()
	pb.RegisterArticleServiceServer(s, srv)

	// listener
	portStr := secret["ARTICLE_SERVICE_PORT"]
	if portStr == nil || portStr == "" {
		log.Fatalf("ARTICLE_SERVICE_PORT is not set")
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
