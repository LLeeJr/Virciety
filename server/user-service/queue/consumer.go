package queue

import (
	"github.com/streadway/amqp"
	"log"
	"user-service/database"
)

const QueryQueue = "user-service-query"
const CommandQueue = "user-service-command"
const EventQueue = "user-service-event"

type Consumer interface {
	InitConsumer(ch *amqp.Channel)
}

func NewConsumer(repo database.Repository) (Consumer, error) {
	return &ChannelConfig{
		Repo: repo,
	}, nil
}

func (c *ChannelConfig) InitConsumer(ch *amqp.Channel) {

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

	if err != nil {
		log.Fatal(err)
	}

	commands, err := ch.Consume(
		CommandQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	events, err := ch.Consume(
		EventQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for data := range queries {
			log.Printf("Received a message with messageID %s : %s", data.MessageId, data.Body)
		}
	}()

	go func() {
		for data := range commands {
			log.Printf("Received a query with messageID %s : %s", data.MessageId, data.Body)
		}
	}()

	go func() {
		for data := range events {
			log.Printf("Received a query with messageID %s : %s", data.MessageId, data.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages")
	<-forever
}