package message_queue

import (
	"event-service/database"
	"github.com/streadway/amqp"
	"log"
)

const QueryQueue = "event-service-query"
const CommandQueue = "event-service-command"
const EventQueue = "event-service-event"

type Consumer interface {
	InitConsumer(ch *amqp.Channel)
}

type ConsumerConfig struct {
	Repo database.Repository
}

func NewConsumer(repo database.Repository) (Consumer, error) {
	return &ConsumerConfig{
		Repo: repo,
	}, nil
}

func (consumer *ConsumerConfig) InitConsumer(ch *amqp.Channel) {
	initQueue(QueryQueue, QueryExchange, ch)
	initQueue(CommandQueue, CommandExchange, ch)
	initQueue(EventQueue, EventExchange, ch)

	queries, _ := ch.Consume(
		QueryQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	commands, _ := ch.Consume(
		CommandQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	events, _ := ch.Consume(
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
			log.Printf("Received a query message with messageID %s : %s", data.MessageId, data.Body)
		}
	}()

	go func() {
		for data := range commands {
			log.Printf("Received a command message with messageID %s : %s", data.MessageId, data.Body)
		}
	}()

	go func() {
		for data := range events {
			log.Printf("Received a event message with messageID %s : %s", data.MessageId, data.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages")
	<-forever
}
