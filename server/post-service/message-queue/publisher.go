package message_queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"post-service/database"
	"post-service/graph/model"
)

const QueryExchange = "query-exchange"
const CommandExchange = "command-exchange"
const EventExchange = "event-exchange"

type Publisher interface {
	InitPublisher(ch *amqp.Channel)
	AddMessageToQuery(postID string, requestID string)
	AddMessageToCommand(comment model.Comment)
	AddMessageToEvent(postEvent database.PostEvent)
	ProfilePictureIdQuery(postId string, requestID string, users []string)
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

func (publisher *PublisherConfig) ProfilePictureIdQuery(postId string, requestID string, users []string) {
	publisher.QueryChan <- RabbitMsg{
		QueueName: QueryExchange,
		PostID:    postId,
		CorrID:    requestID,
		ReplyTo:   QueryExchange,
		Payload:   users,
	}
}

func (publisher *PublisherConfig) AddMessageToQuery(postID string, requestID string) {
	publisher.QueryChan <- RabbitMsg{
		QueueName: QueryExchange,
		PostID:    postID,
		CorrID:    requestID,
		ReplyTo:   QueryExchange,
	}
}

func (publisher *PublisherConfig) AddMessageToCommand(comment model.Comment) {
	publisher.CommandChan <- RabbitMsg{
		QueueName: CommandExchange,
		Comment:   comment,
	}
}

func (publisher *PublisherConfig) AddMessageToEvent(postEvent database.PostEvent) {
	publisher.EventChan <- RabbitMsg{
		QueueName: EventExchange,
		PostEvent: postEvent,
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

type Body struct {
	Id      string      `json:"id"`
	Payload []string    `json:"payload"`
}

func (publisher *PublisherConfig) publish(msg RabbitMsg, ch *amqp.Channel) {
	var body []byte
	var err error
	var corrID = ""
	if msg.QueueName == QueryExchange {
		corrID = msg.CorrID
		if len(msg.Payload) != 0 {
			b := Body{
				Id: msg.PostID,
				Payload: msg.Payload,
			}
			body, err = json.Marshal(b)
		} else {
			b := Body{
				Id: msg.PostID,
			}
			body, err = json.Marshal(b)
		}
	} else if msg.QueueName == CommandExchange {
		body, err = json.Marshal(msg.Comment)
	} else {
		body, err = json.Marshal(msg.PostEvent)
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
			MessageId:     "Post-Service",
			ReplyTo:       msg.ReplyTo,
			CorrelationId: corrID,
		})
	FailOnError(err, "Failed to publish a message")

	log.Printf("INFO: published msg on %s: %v", msg.QueueName, msg.PostEvent)
	log.Printf("ReplyTo: %s, CorrelationID: %s", msg.ReplyTo, msg.CorrID)
}
