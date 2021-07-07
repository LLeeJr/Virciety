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
	InitPublisher()
	AddMessageToQuery()
	AddMessageToCommand(messageId string)
	AddMessageToEvent(dmEvent database.DmEvent, messageId string)
}

func NewPublisher() (Publisher, error) {
	return &ChannelConfig{
		QueryChan:   make(chan RabbitMsg, 10),
		CommandChan: make(chan RabbitMsg, 10),
		EventChan:   make(chan RabbitMsg, 10),
	}, nil
}

func (c *ChannelConfig) InitPublisher() {
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

func (c ChannelConfig) AddMessageToEvent(dmEvent database.DmEvent, messageId string) {
	c.EventChan <- RabbitMsg{
		QueueName: EventExchange,
		DmEvent:   dmEvent,
		MessageId: messageId,
	}
}

func (c *ChannelConfig) publish(msg RabbitMsg, ch *amqp.Channel) {
	body, err := json.Marshal(msg.DmEvent)

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
	failOnError(err, "Failed to publish a message")

	log.Printf(" [*] published msg on %s: %v", msg.QueueName, msg.DmEvent)
}