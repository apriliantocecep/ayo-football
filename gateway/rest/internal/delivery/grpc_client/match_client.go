package grpc_client

import (
	"fmt"
	"github.com/apriliantocecep/ayo-football/services/match/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type MatchServiceClient struct {
	Client pb.MatchServiceClient
	Conn   *grpc.ClientConn
}

func NewMatchServiceClient(vaultClient *shared.VaultClient) *MatchServiceClient {
	secret := utils.GetVaultSecretConfig(vaultClient)

	port := secret["MATCH_SERVICE_PORT"]
	if port == nil || port == "" {
		log.Fatalln("MATCH_SERVICE_PORT is not set")
	}
	url := secret["MATCH_SERVICE_URL"]
	if url == nil || url == "" {
		log.Fatalln("MATCH_SERVICE_URL is not set")
	}

	target := fmt.Sprintf("%s:%s", url.(string), port.(string))
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to match service: %v", err)
	}

	client := pb.NewMatchServiceClient(conn)
	return &MatchServiceClient{
		Client: client,
		Conn:   conn,
	}
}
