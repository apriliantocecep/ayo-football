package main

import (
	"database/sql"
	"fmt"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/config"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/delivery/grpc_server"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/repository"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/usecase"
	"github.com/apriliantocecep/posfin-blog/services/moderation/pkg/pb"
	"github.com/apriliantocecep/posfin-blog/shared"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
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
	metadataRepository := repository.NewMetadataRepository()
	metadataUseCase := usecase.NewMetadataUseCase(database.DB, metadataRepository)

	// grpc server
	srv := grpc_server.NewModerationServer(metadataUseCase)
	s := grpc.NewServer()
	pb.RegisterModerationServiceServer(s, srv)

	// listener
	portStr := secret["MODERATION_SERVICE_PORT"]
	if portStr == nil || portStr == "" {
		log.Fatalln("MODERATION_SERVICE_PORT is not set")
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
