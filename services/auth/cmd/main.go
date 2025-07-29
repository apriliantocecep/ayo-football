package main

import (
	"database/sql"
	"fmt"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/config"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/delivery/grpc_server"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/repository"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/usecase"
	"github.com/apriliantocecep/ayo-football/services/auth/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
)

func main() {
	// setup vars
	serviceName := "auth-service-cluster"

	// vault client
	vaultClient := shared.NewVaultClient()
	secret := utils.GetVaultSecretConfig(vaultClient)

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
	userUseCase := usecase.NewUserUseCase(userRepository, jwt, database.DB)

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

	// health service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
