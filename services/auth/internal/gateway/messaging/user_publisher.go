package messaging

import (
	sharedmessaging "github.com/apriliantocecep/posfin-blog/shared/messaging"
	sharedmodel "github.com/apriliantocecep/posfin-blog/shared/model"
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
