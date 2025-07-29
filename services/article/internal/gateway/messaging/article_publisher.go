package messaging

import (
	sharedmessaging "github.com/apriliantocecep/ayo-football/shared/messaging"
	sharedmodel "github.com/apriliantocecep/ayo-football/shared/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ArticlePublisher struct {
	sharedmessaging.Publisher[*sharedmodel.ArticleEvent]
}

func NewArticlePublisher(rabbitMQConn *amqp.Connection, queueName string, routingKey string) *ArticlePublisher {
	return &ArticlePublisher{
		Publisher: sharedmessaging.Publisher[*sharedmodel.ArticleEvent]{
			RabbitMQConn: rabbitMQConn,
			QueueName:    queueName,
			RoutingKey:   routingKey,
			ExchangeName: "article",
		},
	}
}
