package messaging

import (
	sharedmessaging "github.com/apriliantocecep/ayo-football/shared/messaging"
	sharedmodel "github.com/apriliantocecep/ayo-football/shared/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserPublisher struct {
	sharedmessaging.Publisher[*sharedmodel.UserEvent]
}

func NewUserPublisher(rabbitMQConn *amqp.Connection, queueName string, routingKey string) *UserPublisher {
	return &UserPublisher{
		Publisher: sharedmessaging.Publisher[*sharedmodel.UserEvent]{
			RabbitMQConn: rabbitMQConn,
			QueueName:    queueName,
			RoutingKey:   routingKey,
			ExchangeName: "article",
		},
	}
}
