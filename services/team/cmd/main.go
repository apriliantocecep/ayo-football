package main

import (
	"database/sql"
	"fmt"
	"github.com/apriliantocecep/ayo-football/services/team/internal/config"
	"github.com/apriliantocecep/ayo-football/services/team/internal/delivery/grpc_server"
	"github.com/apriliantocecep/ayo-football/services/team/internal/repository"
	"github.com/apriliantocecep/ayo-football/services/team/internal/usecase"
	"github.com/apriliantocecep/ayo-football/services/team/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
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
	teamRepository := repository.NewTeamRepository()
	teamUseCase := usecase.NewTeamUseCase(database.DB, teamRepository)

	// grpc server
	srv := grpc_server.NewTeamServer(teamUseCase)
	s := grpc.NewServer()
	pb.RegisterTeamServiceServer(s, srv)

	// listener
	portStr := secret["TEAM_SERVICE_PORT"]
	if portStr == nil || portStr == "" {
		log.Fatalln("TEAM_SERVICE_PORT is not set")
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
