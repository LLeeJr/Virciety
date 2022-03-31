package queue

import (
	"github.com/streadway/amqp"
	"log"
	"notifs-service/database"
)

// RabbitMsg and ChannelConfig are both util structs for handling rabbitMQ related data and channels
type (
	RabbitMsg struct {
		QueueName string           `json:"queueName"`
		NotifEvent   database.NotifEvent `json:"dmEvent"`
		MessageId string           `json:"messageId"`
	}

	ChannelConfig struct {
		QueryChan chan RabbitMsg
		CommandChan chan RabbitMsg
		EventChan chan RabbitMsg
		Repo database.Repository
	}
)

// FailOnError is a helper function for logging errors with graceful system exit
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// initExchange declares a queue-exchange for a given queueName and a rabbitMQ-channel
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

// initQueue declares a queue for a given queueName, an rabbitMQ-exchange and a rabbitMQ-channel
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
