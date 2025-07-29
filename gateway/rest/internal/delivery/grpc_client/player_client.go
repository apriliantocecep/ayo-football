package grpc_client

import (
	"fmt"
	"github.com/apriliantocecep/ayo-football/services/player/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type PlayerServiceClient struct {
	Client pb.PlayerServiceClient
	Conn   *grpc.ClientConn
}

func NewPlayerServiceClient(vaultClient *shared.VaultClient) *PlayerServiceClient {
	secret := utils.GetVaultSecretConfig(vaultClient)

	port := secret["PLAYER_SERVICE_PORT"]
	if port == nil || port == "" {
		log.Fatalln("PLAYER_SERVICE_PORT is not set")
	}
	url := secret["PLAYER_SERVICE_URL"]
	if url == nil || url == "" {
		log.Fatalln("PLAYER_SERVICE_URL is not set")
	}

	target := fmt.Sprintf("%s:%s", url.(string), port.(string))
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to team service: %v", err)
	}

	client := pb.NewPlayerServiceClient(conn)
	return &PlayerServiceClient{
		Client: client,
		Conn:   conn,
	}
}
