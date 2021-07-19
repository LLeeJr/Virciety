package message_queue

import (
	"comment-service/database"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMsg struct {
	QueueName    string                `json:"queueName"`
	CommentEvent database.CommentEvent `json:"commentEvent"`
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

func initQueue(queueName string, exchangeName string, ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		exchangeName,
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
