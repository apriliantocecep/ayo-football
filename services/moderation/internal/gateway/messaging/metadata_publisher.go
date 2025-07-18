package messaging

import (
	sharedmessaging "github.com/apriliantocecep/posfin-blog/shared/messaging"
	sharedmodel "github.com/apriliantocecep/posfin-blog/shared/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MetadataPublisher struct {
	sharedmessaging.Publisher[*sharedmodel.MetadataEvent]
}

func NewMetadataPublisher(channel *amqp.Channel, queueName string, routingKey string) *MetadataPublisher {
	return &MetadataPublisher{
		Publisher: sharedmessaging.Publisher[*sharedmodel.MetadataEvent]{
			Channel:      channel,
			QueueName:    queueName,
			RoutingKey:   routingKey,
			ExchangeName: "article",
		},
	}
}
