package message_queue

import (
	"github.com/streadway/amqp"
	"log"
	"post-service/database"
	"post-service/graph/model"
)

type RabbitMsg struct {
	QueueName string             `json:"queueName"`
	PostEvent database.PostEvent `json:"postEvent"`
	Comment   model.Comment      `json:"comment"`
	PostID    string             `json:"postID"`
	CorrID    string             `json:"corrID"`
	ReplyTo   string			 `json:"replyTo"`
	Payload   []string           `json:"payload"`
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
	FailOnError(err, "Failed to declare an exchange")
}

func FailOnError(err error, msg string) {
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
	FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		exchangeName,
		false,
		nil)
	FailOnError(err, "Failed to bind a queue")
}
