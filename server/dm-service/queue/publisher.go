package queue

import (
	"dm-service/database"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

const QueryExchange = "query-exchange"
const CommandExchange = "command-exchange"
const EventExchange = "event-exchange"

// Publisher is a helper-interface for initialising a rabbitMQ-publisher and handling its message-production
type Publisher interface {
	InitPublisher(ch *amqp.Channel)
	AddMessageToQuery()
	AddMessageToCommand(messageId string)
	AddMessageToEvent(chatEvent database.ChatEvent, messageId string)
}

// NewPublisher creates a new Publisher with newly created channels
func NewPublisher() (Publisher, error) {
	return &ChannelConfig{
		QueryChan:   make(chan RabbitMsg, 10),
		CommandChan: make(chan RabbitMsg, 10),
		EventChan:   make(chan RabbitMsg, 10),
	}, nil
}

// InitPublisher initialises a Publisher and its rabbitMQ-exchanges
func (c *ChannelConfig) InitPublisher(ch *amqp.Channel) {

	initExchange(QueryExchange, ch)
	initExchange(CommandExchange, ch)
	initExchange(EventExchange, ch)

	for {
		select {
		case msg := <-c.QueryChan:
			c.publish(msg, ch)
		case msg := <-c.CommandChan:
			c.publish(msg, ch)
		case msg := <-c.EventChan:
			c.publish(msg, ch)
		}
	}
}

// AddMessageToQuery adds a message to the query-exchange
func (c ChannelConfig) AddMessageToQuery() {
	c.QueryChan <- RabbitMsg{
		QueueName: QueryExchange,
	}
}

// AddMessageToCommand adds a message to the command-exchange
func (c ChannelConfig) AddMessageToCommand(messageId string) {
	c.CommandChan <- RabbitMsg{
		QueueName: CommandExchange,
		MessageId: messageId,
	}
}

// AddMessageToEvent adds a message to the event-exchange
func (c ChannelConfig) AddMessageToEvent(chatEvent database.ChatEvent, messageId string) {
	c.EventChan <- RabbitMsg{
		QueueName: EventExchange,
		ChatEvent: chatEvent,
		MessageId: messageId,
	}
}

// publish is a helper-function for producing new messages on rabbitMQ
func (c *ChannelConfig) publish(msg RabbitMsg, ch *amqp.Channel) {
	body, err := json.Marshal(msg.ChatEvent)

	err = ch.Publish(
		msg.QueueName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			MessageId:   msg.MessageId,
		})
	FailOnError(err, "Failed to publish a message")

	log.Printf(" [*] published msg on %s: %v", msg.QueueName, msg.ChatEvent)
}