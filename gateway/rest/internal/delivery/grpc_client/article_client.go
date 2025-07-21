package grpc_client

import (
	"github.com/apriliantocecep/posfin-blog/services/article/pkg/pb"
	"github.com/apriliantocecep/posfin-blog/shared"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ArticleServiceClient struct {
	Client pb.ArticleServiceClient
	Conn   *grpc.ClientConn
}

func NewArticleServiceClient(vaultClient *shared.VaultClient) *ArticleServiceClient {
	secret := utils.GetVaultSecretConfig(vaultClient)

	proxyUrl := secret["ARTICLE_SERVICE_PROXY"]
	if proxyUrl == nil || proxyUrl == "" {
		log.Fatalln("ARTICLE_SERVICE_PROXY is not set")
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
		log.Fatalf("did not connect to article service: %v", err)
	}

	client := pb.NewArticleServiceClient(conn)

	return &ArticleServiceClient{
		Client: client,
		Conn:   conn,
	}
}
