package message_queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"post-service/database"
	"post-service/graph/model"
)

const QueryQueue = "post-service-query"
const CommandQueue = "post-service-command"
const EventQueue = "post-service-event"

type Consumer interface {
	InitConsumer(ch *amqp.Channel)
}

type ConsumerConfig struct {
	Repo          database.Repository
	Responses     map[string]chan []*model.Comment
	UserResponses map[string]chan map[string]string
}

func NewConsumer(repo database.Repository, responses map[string]chan []*model.Comment, userResponses map[string]chan map[string]string) (Consumer, error) {
	return &ConsumerConfig{
		Repo: 		   repo,
		Responses: 	   responses,
		UserResponses: userResponses,
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
			if data.MessageId == "Comment-Service" {
				var comments []*model.Comment
				err := json.Unmarshal(data.Body, &comments)
				if err != nil {
					log.Fatalln(err)
				}

				consumer.Responses[data.CorrelationId] <- comments
			}
			if data.MessageId == "User-Service" {
				var userIdMap map[string]string
				err := json.Unmarshal(data.Body, &userIdMap)
				if err != nil {
					log.Fatalln(err)
				}
				consumer.UserResponses[data.CorrelationId] <- userIdMap
			}
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
