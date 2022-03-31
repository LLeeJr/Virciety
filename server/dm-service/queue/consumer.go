package queue

import (
	"dm-service/database"
	"github.com/streadway/amqp"
	"log"
)

const QueryQueue = "dm-service-query"
const CommandQueue = "dm-service-command"
const EventQueue = "dm-service-event"

// Consumer is a helper-interface for initialising a rabbitMQ-consumer and handling its message-consumption
type Consumer interface {
	InitConsumer(ch *amqp.Channel)
}

// NewConsumer creates a new Consumer with a respective repo to make calls to the database
func NewConsumer(repo database.Repository) (Consumer, error) {
	return &ChannelConfig{
		Repo: repo,
	}, nil
}

// InitConsumer initialises a Consumer and its rabbitMQ-exchanges, as well as activates its queue-listeners
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