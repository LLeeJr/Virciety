package message_queue

import (
	"comment-service/database"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

const QueryExchange = "query-exchange"
const CommandExchange = "command-exchange"
const EventExchange = "event-exchange"

type Publisher interface {
	InitPublisher()
	AddMessageToQuery()
	AddMessageToCommand()
	AddMessageToEvent(commentEvent database.CommentEvent)
}

func NewPublisher() (Publisher, error) {
	return &ChannelConfig{
		QueryChan:   make(chan RabbitMsg, 10),
		CommandChan: make(chan RabbitMsg, 10),
		EventChan:   make(chan RabbitMsg, 10),
	}, nil
}

func (channel *ChannelConfig) AddMessageToQuery() {
	channel.QueryChan <- RabbitMsg{
		QueueName: QueryExchange,
	}
}

func (channel *ChannelConfig) AddMessageToCommand() {
	channel.CommandChan <- RabbitMsg{
		QueueName: CommandExchange,
	}
}

func (channel *ChannelConfig) AddMessageToEvent(commentEvent database.CommentEvent) {
	channel.EventChan <- RabbitMsg{
		QueueName:    EventExchange,
		CommentEvent: commentEvent,
	}
}

func (channel *ChannelConfig) InitPublisher() {
	// conn
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	initExchange(QueryExchange, ch)
	initExchange(CommandExchange, ch)
	initExchange(EventExchange, ch)

	for {
		select {
		case msg := <-channel.QueryChan:
			channel.publish(msg, ch)
		case msg := <-channel.CommandChan:
			channel.publish(msg, ch)
		case msg := <-channel.EventChan:
			channel.publish(msg, ch)
		}
	}
}

func (channel *ChannelConfig) publish(msg RabbitMsg, ch *amqp.Channel) {
	body, err := json.Marshal(msg.CommentEvent)

	// publish message
	err = ch.Publish(
		msg.QueueName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
			MessageId:   "Comment-Service",
		})
	failOnError(err, "Failed to publish a message")

	log.Printf("INFO: published msg on %s: %v", msg.QueueName, msg.CommentEvent)
}
