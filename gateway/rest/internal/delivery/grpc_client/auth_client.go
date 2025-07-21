package grpc_client

import (
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

	proxyUrl := secret["AUTH_SERVICE_PROXY"]
	if proxyUrl == nil || proxyUrl == "" {
		log.Fatalln("AUTH_SERVICE_PROXY is not set")
	}
	grpcProxyUrl := secret["TRAEFIK_GRPC_PROXY_URL"]
	if grpcProxyUrl == nil || grpcProxyUrl == "" {
		log.Fatalln("TRAEFIK_GRPC_PROXY_URL is not set")
	}

	conn, err := grpc.NewClient(
		grpcProxyUrl.(string),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithAuthority(proxyUrl.(string)),
	)
	if err != nil {
		log.Fatalf("did not connect to auth service: %v", err)
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthServiceClient{
		Client: client,
		Conn:   conn,
	}
}
