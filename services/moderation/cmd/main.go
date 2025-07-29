package main

import (
	"database/sql"
	"fmt"
	"github.com/apriliantocecep/ayo-football/services/moderation/internal/config"
	"github.com/apriliantocecep/ayo-football/services/moderation/internal/delivery/grpc_server"
	"github.com/apriliantocecep/ayo-football/services/moderation/internal/repository"
	"github.com/apriliantocecep/ayo-football/services/moderation/internal/usecase"
	"github.com/apriliantocecep/ayo-football/services/moderation/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared"
	sharedlib "github.com/apriliantocecep/ayo-football/shared/lib"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	capi "github.com/hashicorp/consul/api"
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

	// grpc server
	srv := grpc_server.NewModerationServer(metadataUseCase)
	s := grpc.NewServer()
	pb.RegisterModerationServiceServer(s, srv)

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
	serviceName := "moderation-service-cluster"
	serviceRegisteredID := fmt.Sprintf("moderation-service-%d", port)
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
