package queue

import (
	"github.com/streadway/amqp"
	"log"
	"notifs-service/database"
)

const QueryQueue = "notif-service-query"
const CommandQueue = "notif-service-command"
const EventQueue = "notif-service-event"

type Consumer interface {
	InitConsumer()
}

func NewConsumer(repo database.Repository) (Consumer, error) {
	return &ChannelConfig{
		Repo: repo,
	}, nil
}

func (c *ChannelConfig) InitConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	initQueue(QueryQueue, QueryExchange, ch)
	initQueue(CommandQueue, CommandExchange, ch)
	initQueue(EventQueue, EventExchange, ch)

	queries, err := ch.Consume(
		QueryQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	commands, err := ch.Consume(
		CommandQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	events, err := ch.Consume(
		EventQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for data := range queries {
			log.Printf("Query: Received a message with messageID %s : %s", data.MessageId, data.Body)
		}
	}()

	go func() {
		for data := range commands {
			log.Printf("Command: Received a query with messageID %s : %s", data.MessageId, data.Body)
		}
	}()

	go func() {
		for data := range events {
			log.Printf("Event: Received a query with messageID %s : %s", data.MessageId, data.Body)
			if data.MessageId == "Dm-Service" {
				c.Repo.CreateDmNotifFromConsumer(data.Body)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages")
	<-forever
}