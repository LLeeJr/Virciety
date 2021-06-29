package message_queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"posts-service/database"
)

type MessageQueue interface {
	InitProducer()
	AddMessageToQueue(postEvent database.PostEvent)
	AddMessageToExchange(postEvent database.PostEvent)
}

type RabbitMsg struct {
	QueueName string             `json:"queueName"`
	PostEvent database.PostEvent `json:"postEvent"`
}

type Channel struct {
	QChan chan RabbitMsg
	EChan chan RabbitMsg
}

func NewChannel() (MessageQueue, error) {
	return &Channel{QChan: make(chan RabbitMsg, 10),
		EChan: make(chan RabbitMsg, 10)}, nil
}

func (channel *Channel) AddMessageToQueue(postEvent database.PostEvent) {
	channel.QChan <- RabbitMsg{
		QueueName: "notification-service",
		PostEvent: postEvent,
	}
}

func (channel *Channel) AddMessageToExchange(postEvent database.PostEvent) {
	channel.EChan <- RabbitMsg{
		QueueName: "",
		PostEvent: postEvent,
	}
}

func (channel *Channel) InitProducer() {
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

			log.Printf("INFO: published msg: %v", msg.PostEvent)
		case msg := <-channel.EChan:
			log.Printf("INFO: got msg: %v", msg.PostEvent)
		}
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
