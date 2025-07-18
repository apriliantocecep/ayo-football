package messaging

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type ConsumerHandler func(delivery amqp.Delivery) error

func ConsumeQueue(ctx context.Context, channel *amqp.Channel, queueName string, handler ConsumerHandler) {
	// NOTE: Because we might start the consumer before the publisher, we want to make sure
	// the queue exists before we try to consume messages from it.
	q, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-queue-type": "quorum",
		}, // arguments
	)
	if err != nil {
		log.Fatalf("consumer failed to declare a queue '%s' : %v", queueName, err)
	}

	msgs, err := channel.ConsumeWithContext(ctx,
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("consumer failed register to queue '%s': %v", queueName, err)
	}

	run := true

	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			for d := range msgs {
				err = handler(d)
				if err != nil {
					log.Fatalf("failed to process message: %v", err)
				}
			}
		}
	}

	log.Printf("closing consumer for queue: %s", queueName)
}
