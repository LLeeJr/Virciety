package message_queue

import (
	"github.com/streadway/amqp"
	"log"
	"posts-service/database"
	"posts-service/graph/model"
)

type RabbitMsg struct {
	QueueName string             `json:"queueName"`
	PostEvent database.PostEvent `json:"postEvent"`
	Comment   model.Comment      `json:"comment"`
}

type ChannelConfig struct {
	QueryChan   chan RabbitMsg
	CommandChan chan RabbitMsg
	EventChan   chan RabbitMsg
	Repo        database.Repository
}

type Comment struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	LikedBy     []string `json:"likedBy"`
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
