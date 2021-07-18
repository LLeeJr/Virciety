package message_queue

import (
	"comment-service/database"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMsg struct {
	QueueName    string                `json:"queueName"`
	CommentEvent database.CommentEvent `json:"commentEvent"`
	MessageId    string                `json:"messageId"`
}

type ChannelConfig struct {
	QueryChan   chan RabbitMsg
	CommandChan chan RabbitMsg
	EventChan   chan RabbitMsg
	Repo        database.Repository
}

func initExchange(queueName string, ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		queueName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
