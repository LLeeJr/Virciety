package queue

import (
	"github.com/streadway/amqp"
	"log"
	"notifs-service/database"
)

const QueryQueue = "notif-service-query"
const CommandQueue = "notif-service-command"
const EventQueue = "notif-service-event"

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
			if data.MessageId == "Comment-Service" {
				c.Repo.CreateCommentNotifFromConsumer(data.Body)
			}
			if data.MessageId == "Post-Service" {
				c.Repo.CreateLikeNotifFromConsumer(data.Body)
			}
			if data.MessageId == "User-Service" {
				c.Repo.CreateFollowNotifFromConsumer(data.Body)
			}
			if data.MessageId == "Event-Service" {
				c.Repo.CreateEventNotifFromConsumer(data.Body)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages")
	<-forever
}