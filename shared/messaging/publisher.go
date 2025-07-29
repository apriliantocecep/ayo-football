package messaging

import (
	"context"
	"encoding/json"
	sharedmodel "github.com/apriliantocecep/ayo-football/shared/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type Publisher[T sharedmodel.Event] struct {
	RabbitMQConn *amqp.Connection
	QueueName    string
	RoutingKey   string
	ExchangeName string
}

func (p *Publisher[T]) GetQueueName() string {
	return p.QueueName
}

func (p *Publisher[T]) GetRoutingKey() string {
	return p.RoutingKey
}

func (p *Publisher[T]) GetExchangeName() string {
	return p.ExchangeName
}

func (p *Publisher[T]) Publish(event T) error {
	value, err := json.Marshal(event)
	if err != nil {
		log.Println("failed to marshal event")
	}

	channel, err := p.RabbitMQConn.Channel()
	if err != nil {
		log.Printf("[publisher] failed to open a channel for queue '%s' : %v", p.GetQueueName(), err)
		return err
	}
	defer func(ch *amqp.Channel) {
		err = ch.Close()
		if err != nil {
			log.Fatalf("[publisher] failed closing a channel for queue '%s' : %v", p.GetQueueName(), err)
		}
	}(channel)

	_, err = channel.QueueDeclare(
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

	err = channel.PublishWithContext(ctx,
		p.GetExchangeName(), // exchange
		p.GetRoutingKey(),   // routing key
		false,               // mandatory
		false,               // immediate
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
