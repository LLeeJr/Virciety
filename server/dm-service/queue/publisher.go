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

type Publisher interface {
	InitPublisher(ch *amqp.Channel)
	AddMessageToQuery()
	AddMessageToCommand(messageId string)
	AddMessageToEvent(chatEvent database.ChatEvent, messageId string)
}

func NewPublisher() (Publisher, error) {
	return &ChannelConfig{
		QueryChan:   make(chan RabbitMsg, 10),
		CommandChan: make(chan RabbitMsg, 10),
		EventChan:   make(chan RabbitMsg, 10),
	}, nil
}

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

func (c ChannelConfig) AddMessageToQuery() {
	c.QueryChan <- RabbitMsg{
		QueueName: QueryExchange,
	}
}

func (c ChannelConfig) AddMessageToCommand(messageId string) {
	c.CommandChan <- RabbitMsg{
		QueueName: CommandExchange,
		MessageId: messageId,
	}
}

func (c ChannelConfig) AddMessageToEvent(chatEvent database.ChatEvent, messageId string) {
	c.EventChan <- RabbitMsg{
		QueueName: EventExchange,
		ChatEvent: chatEvent,
		MessageId: messageId,
	}
}

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