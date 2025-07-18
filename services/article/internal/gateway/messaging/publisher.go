package messaging

import (
	"context"
	"encoding/json"
	sharedmodel "github.com/apriliantocecep/posfin-blog/shared/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type Publisher[T sharedmodel.Event] struct {
	Channel    *amqp.Channel
	QueueName  string
	RoutingKey string
}

func (p *Publisher[T]) GetQueueName() string {
	return p.QueueName
}

func (p *Publisher[T]) GetRoutingKey() string {
	return p.RoutingKey
}

func (p *Publisher[T]) Publish(event T) error {
	value, err := json.Marshal(event)
	if err != nil {
		log.Println("failed to marshal event")
	}

	_, err = p.Channel.QueueDeclare(
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
		log.Printf("failed to declare a queue '%s' : %v", p.GetQueueName(), err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.Channel.PublishWithContext(ctx,
		"article",         // exchange
		p.GetRoutingKey(), // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        value,
			Timestamp:   time.Now(),
		})
	if err != nil {
		log.Printf("failed publish a message to '%s': %v", p.GetRoutingKey(), err)
		return err
	}

	return nil
}
