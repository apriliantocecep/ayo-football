package grpc_client

import (
	"fmt"
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

	port := secret["ARTICLE_SERVICE_PORT"]
	if port == nil || port == "" {
		log.Fatalln("ARTICLE_SERVICE_PORT is not set")
	}
	url := secret["ARTICLE_SERVICE_URL"]
	if url == nil || url == "" {
		log.Fatalln("ARTICLE_SERVICE_URL is not set")
	}

	target := fmt.Sprintf("%s:%s", url.(string), port.(string))
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to auth service: %v", err)
	}

	client := pb.NewArticleServiceClient(conn)

	return &ArticleServiceClient{
		Client: client,
		Conn:   conn,
	}
}
