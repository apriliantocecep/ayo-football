package main

import (
	"context"
	"fmt"
	"github.com/apriliantocecep/ayo-football/services/article/internal/config"
	"github.com/apriliantocecep/ayo-football/services/article/internal/delivery/grpc_server"
	"github.com/apriliantocecep/ayo-football/services/article/internal/gateway/messaging"
	"github.com/apriliantocecep/ayo-football/services/article/internal/repository"
	"github.com/apriliantocecep/ayo-football/services/article/internal/usecase"
	"github.com/apriliantocecep/ayo-football/services/article/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared"
	sharedlib "github.com/apriliantocecep/ayo-football/shared/lib"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	capi "github.com/hashicorp/consul/api"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
)

func main() {
	// vault client
	vaultClient := shared.NewVaultClient()
	//secret := utils.GetVaultSecretConfig(vaultClient)

	// consul client
	consul := sharedlib.NewConsulClient(vaultClient)

	// rabbitmq client
	rabbitMQClient := sharedlib.NewRabbitMQClient(vaultClient)
	defer func(Conn *amqp.Connection) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("error closing rabbitmq: %v", err)
		}
	}(rabbitMQClient.Conn)

	// setup publisher
	articleCreatedPublisher := messaging.NewArticlePublisher(rabbitMQClient.Conn, "article_created", "article_created")
	articleModerationPublisher := messaging.NewArticlePublisher(rabbitMQClient.Conn, "article_moderation", "article_moderation")

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
	articleUseCase := usecase.NewArticleUseCase(database.Client, articleRepository, articleCreatedPublisher, articleModerationPublisher)

	// grpc server
	srv := grpc_server.NewArticleServer(articleUseCase)
	s := grpc.NewServer()
	pb.RegisterArticleServiceServer(s, srv)

	// listener
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatalln("PORT is not set")
	}
	port := utils.ParsePort(portStr)

	url := os.Getenv("SERVICE_URL")
	if url == "" {
		log.Fatalln("SERVICE_URL is not set")
	}
	address := url

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", listen.Addr())

	// register service to consul
	serviceName := "article-service-cluster"
	serviceRegisteredID := fmt.Sprintf("article-service-%d", port)
	tags := []string{
		"traefik.enable=true",
		fmt.Sprintf("traefik.http.routers.%s.rule=Host(`%s.local`)", serviceName, serviceName),
		fmt.Sprintf("traefik.http.routers.%s.entrypoints=grpc", serviceName),
		fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.scheme=h2c", serviceName),
	}
	consul.ServiceRegister(serviceRegisteredID, serviceName, address, port, &capi.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", address, port),
		Interval:                       "10s",
		Timeout:                        "2s",
		DeregisterCriticalServiceAfter: "1m",
	}, tags)

	// health service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
