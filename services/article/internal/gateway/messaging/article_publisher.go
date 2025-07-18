package messaging

import (
	sharedmodel "github.com/apriliantocecep/posfin-blog/shared/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ArticlePublisher struct {
	Publisher[*sharedmodel.ArticleEvent]
}

func NewArticlePublisher(channel *amqp.Channel, queueName string, routingKey string) *ArticlePublisher {
	return &ArticlePublisher{
		Publisher: Publisher[*sharedmodel.ArticleEvent]{
			Channel:    channel,
			QueueName:  queueName,
			RoutingKey: routingKey,
		},
	}
}
