package queue

import (
	"dm-service/database"
	"github.com/streadway/amqp"
	"log"
)

type (
	RabbitMsg struct {
		QueueName string           `json:"queueName"`
		DmEvent   database.DmEvent `json:"dmEvent"`
		MessageId string           `json:"messageId"`
	}

	ChannelConfig struct {
		QueryChan chan RabbitMsg
		CommandChan chan RabbitMsg
		EventChan chan RabbitMsg
		Repo database.Repository
	}
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func initExchange(queueName string, ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		queueName,
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	FailOnError(err, "Failed to declare an exchange")
}

func initQueue(queueName string, exchangeName string, ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		queueName,    // name
		true, // durable
		false, // delete when unused
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		exchangeName,
		false,
		nil,
	)
	FailOnError(err, "Failed to bind a queue")
}
