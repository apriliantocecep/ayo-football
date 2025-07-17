package grpc_client

import (
	"fmt"
	"github.com/apriliantocecep/posfin-blog/services/auth/pkg/pb"
	"github.com/apriliantocecep/posfin-blog/shared"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
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

	port := secret["AUTH_SERVICE_PORT"].(string)
	url := secret["AUTH_SERVICE_URL"].(string)

	target := fmt.Sprintf("%s:%s", url, port)
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
