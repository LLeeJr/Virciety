package message_queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"posts-service/database"
)

type ProducerQueue interface {
	InitProducer()
	AddMessageToQueue(postEvent database.PostEvent)
	AddMessageToExchange(postEvent database.PostEvent)
}

func NewProducerChannel() (ProducerQueue, error) {
	return &ChannelConfig{QChan: make(chan RabbitMsg, 10),
		EChan: make(chan RabbitMsg, 10)}, nil
}

func (channel *ChannelConfig) AddMessageToQueue(postEvent database.PostEvent) {
	channel.QChan <- RabbitMsg{
		QueueName: "notification-service",
		PostEvent: postEvent,
	}
}

func (channel *ChannelConfig) AddMessageToExchange(postEvent database.PostEvent) {
	channel.EChan <- RabbitMsg{
		QueueName: "",
		PostEvent: postEvent,
	}
}

func (channel *ChannelConfig) InitProducer() {
	// conn
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"notification-service",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	for {
		select {
		case msg := <-channel.QChan:
			body, err := json.Marshal(msg.PostEvent)

			// publish message
			err = ch.Publish(
				"",
				msg.QueueName,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        body,
				})
			failOnError(err, "Failed to publish a message")

			log.Printf("INFO: published msg on %s: %v", msg.QueueName, msg.PostEvent)
		case msg := <-channel.EChan:
			log.Printf("INFO: got msg on %s: %v", msg.QueueName, msg.PostEvent)
		}
	}
}
