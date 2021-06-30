package message_queue

import (
	"github.com/streadway/amqp"
	"log"
	"posts-service/database"
)

type ConsumerQueue interface {
	InitConsumer()
}

func NewConsumerChannel(repo database.Repository) (ConsumerQueue, error) {
	return &ChannelConfig{
		Repo: repo,
	}, nil
}

func (channel *ChannelConfig) InitConsumer() {
	// conn
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"comment-service",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for data := range msgs {
			log.Printf("Received a message: %s", data.Body)
			// TODO do stuff with event
		}
	}()

	log.Printf(" [*] Waiting for messages")
	<-forever
}
