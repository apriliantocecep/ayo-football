package grpc_client

import (
	"fmt"
	"github.com/apriliantocecep/ayo-football/services/team/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type TeamServiceClient struct {
	Client pb.TeamServiceClient
	Conn   *grpc.ClientConn
}

func NewTeamServiceClient(vaultClient *shared.VaultClient) *TeamServiceClient {
	secret := utils.GetVaultSecretConfig(vaultClient)

	port := secret["TEAM_SERVICE_PORT"]
	if port == nil || port == "" {
		log.Fatalln("TEAM_SERVICE_PORT is not set")
	}
	url := secret["TEAM_SERVICE_URL"]
	if url == nil || url == "" {
		log.Fatalln("TEAM_SERVICE_URL is not set")
	}

	target := fmt.Sprintf("%s:%s", url.(string), port.(string))
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to team service: %v", err)
	}

	client := pb.NewTeamServiceClient(conn)
	return &TeamServiceClient{
		Client: client,
		Conn:   conn,
	}
}
