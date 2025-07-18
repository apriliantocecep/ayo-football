package messaging

import (
	"github.com/apriliantocecep/posfin-blog/services/article/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ArticlePublisher struct {
	Publisher[*model.ArticleEvent]
}

func NewArticlePublisher(channel *amqp.Channel, queueName string, routingKey string) *ArticlePublisher {
	return &ArticlePublisher{
		Publisher: Publisher[*model.ArticleEvent]{
			Channel:    channel,
			QueueName:  queueName,
			RoutingKey: routingKey,
		},
	}
}
