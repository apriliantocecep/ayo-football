package grpc_client

import (
	"fmt"
	"github.com/apriliantocecep/ayo-football/services/auth/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type AuthServiceClient struct {
	Client pb.AuthServiceClient
	Conn   *grpc.ClientConn
}

func NewAuthServiceClient(vaultClient *shared.VaultClient) *AuthServiceClient {
	secret := utils.GetVaultSecretConfig(vaultClient)

	port := secret["AUTH_SERVICE_PORT"]
	if port == nil || port == "" {
		log.Fatalln("AUTH_SERVICE_PORT is not set")
	}
	url := secret["AUTH_SERVICE_URL"]
	if url == nil || url == "" {
		log.Fatalln("AUTH_SERVICE_URL is not set")
	}

	target := fmt.Sprintf("%s:%s", url.(string), port.(string))
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to auth service: %v", err)
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthServiceClient{
		Client: client,
		Conn:   conn,
	}
}
