package lib

import (
	"github.com/apriliantocecep/posfin-blog/shared"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMQClient struct {
	Conn *amqp.Connection
}

func NewRabbitMQClient(vaultClient *shared.VaultClient) *RabbitMQClient {
	secret := utils.GetVaultSecretConfig(vaultClient)
	url := secret["RABBITMQ_URL"]
	if url == nil || url == "" {
		log.Fatalln("RABBITMQ_URL is not set")
	}

	conn, err := amqp.Dial(url.(string))
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	return &RabbitMQClient{
		Conn: conn,
	}
}
