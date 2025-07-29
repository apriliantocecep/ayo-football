package config

import (
	"context"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"log"
	"time"
)

type Database struct {
	Client *mongo.Client
}

func NewDatabase(vaultClient *shared.VaultClient) *Database {
	secret := utils.GetVaultSecretArticleSvc(vaultClient)

	uri := secret["DATABASE_URL"]
	if uri == nil || uri == "" {
		log.Fatalln("DATABASE_URL is not set")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri.(string)))
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return &Database{Client: client}
}
