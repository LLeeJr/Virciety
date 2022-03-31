package message_queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

const QueryExchange = "query-exchange"
const CommandExchange = "command-exchange"
const EventExchange = "event-exchange"

type EventNotification struct {
	EventId    string `json:"eventID"`
	Message    string `json:"message"`
	ReportedBy string `json:"reportedBy"`
	Username   string `json:"username"`
	EditFlag   bool   `json:"edit_flag"`
}

type Publisher interface {
	InitPublisher(ch *amqp.Channel)
	AddMessageToQuery(postID string, requestID string)
	AddMessageToCommand()
	AddMessageToEvent(eventNotification EventNotification)
}

type PublisherConfig struct {
	QueryChan   chan RabbitMsg
	CommandChan chan RabbitMsg
	EventChan   chan RabbitMsg
}

func NewPublisher() (Publisher, error) {
	return &PublisherConfig{
		QueryChan:   make(chan RabbitMsg, 10),
		CommandChan: make(chan RabbitMsg, 10),
		EventChan:   make(chan RabbitMsg, 10),
	}, nil
}

func (publisher *PublisherConfig) AddMessageToQuery(postID string, requestID string) {
	publisher.QueryChan <- RabbitMsg{
		QueueName: QueryExchange,
		PostID:    postID,
		CorrID:    requestID,
		ReplyTo:   QueryExchange,
	}
}

func (publisher *PublisherConfig) AddMessageToCommand() {
	publisher.CommandChan <- RabbitMsg{
		QueueName: CommandExchange,
	}
}

func (publisher *PublisherConfig) AddMessageToEvent(eventNotification EventNotification) {
	publisher.EventChan <- RabbitMsg{
		QueueName:         EventExchange,
		EventNotification: eventNotification,
	}
}

func (publisher *PublisherConfig) InitPublisher(ch *amqp.Channel) {
	initExchange(QueryExchange, ch)
	initExchange(CommandExchange, ch)
	initExchange(EventExchange, ch)

	for {
		select {
		case msg := <-publisher.QueryChan:
			publisher.publish(msg, ch)
		case msg := <-publisher.CommandChan:
			publisher.publish(msg, ch)
		case msg := <-publisher.EventChan:
			publisher.publish(msg, ch)
		}
	}
}

func (publisher *PublisherConfig) publish(msg RabbitMsg, ch *amqp.Channel) {
	var body []byte
	var err error
	var corrID = ""
	if msg.QueueName == QueryExchange {
		corrID = msg.CorrID
		body, err = json.Marshal(msg.PostID)
	} else if msg.QueueName == CommandExchange {
	} else if msg.QueueName == EventExchange {
		body, err = json.Marshal(msg.EventNotification)
	}
	FailOnError(err, "Failed to json.marshal request")

	// publish message
	err = ch.Publish(
		msg.QueueName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          body,
			MessageId:     "Event-Service",
			ReplyTo:       msg.ReplyTo,
			CorrelationId: corrID,
		})
	FailOnError(err, "Failed to publish a message")

	log.Printf("INFO: published msg on %s: %v", msg.QueueName, msg.EventEvent)
	log.Printf("ReplyTo: %s, CorrelationID: %s", msg.ReplyTo, msg.CorrID)
}
