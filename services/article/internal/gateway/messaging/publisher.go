package messaging

import (
	"context"
	"encoding/json"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type Publisher[T model.Event] struct {
	Channel   *amqp.Channel
	QueueName string
}

func (p *Publisher[T]) GetQueueName() string {
	return p.QueueName
}

func (p *Publisher[T]) Publish(event T) error {
	value, err := json.Marshal(event)
	if err != nil {
		log.Println("failed to marshal event")
	}

	q, err := p.Channel.QueueDeclare(
		p.GetQueueName(), // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		amqp.Table{
			"x-queue-type": "quorum",
		}, // arguments
	)
	if err != nil {
		log.Printf("failed to declare a queue: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.Channel.PublishWithContext(ctx,
		"article", // exchange
		q.Name,    // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        value,
			Timestamp:   time.Now(),
		})
	if err != nil {
		log.Printf("failed to publish a message: %v", err)
	}

	return nil
}
