package message_queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"posts-service/database"
	"posts-service/graph/model"
)

const QueryExchange = "query-exchange"
const CommandExchange = "command-exchange"
const EventExchange = "event-exchange"

type Publisher interface {
	InitPublisher()
	AddMessageToQuery()
	AddMessageToCommand(comment model.Comment)
	AddMessageToEvent(postEvent database.PostEvent)
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

func (channel *ChannelConfig) AddMessageToCommand(comment model.Comment) {
	channel.CommandChan <- RabbitMsg{
		QueueName: CommandExchange,
		Comment:   comment,
	}
}

func (channel *ChannelConfig) AddMessageToEvent(postEvent database.PostEvent) {
	channel.EventChan <- RabbitMsg{
		QueueName: EventExchange,
		PostEvent: postEvent,
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
	var body []byte
	var err error
	//var corrID = uuid.NewString()
	if msg.QueueName == QueryExchange {

	} else if msg.QueueName == CommandExchange {
		body, err = json.Marshal(msg.Comment)
	} else {
		body, err = json.Marshal(msg.PostEvent)
	}
	failOnError(err, "Failed to json.marshal request")

	// publish message
	err = ch.Publish(
		msg.QueueName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
			MessageId:   "Post-Service",
			/*ReplyTo: 	 	msg.ReplyTo,
			CorrelationId: 	corrID,*/
		})
	failOnError(err, "Failed to publish a message")

	log.Printf("INFO: published msg on %s: %v", msg.QueueName, msg.PostEvent)
}
