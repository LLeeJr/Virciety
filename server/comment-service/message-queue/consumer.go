package message_queue

import (
	"comment-service/database"
	"comment-service/model"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

const QueryQueue = "comment-service-query"
const CommandQueue = "comment-service-command"
const EventQueue = "comment-service-event"

type Consumer interface {
	InitConsumer(ch *amqp.Channel)
}

func NewConsumer(repo database.Repository) (Consumer, error) {
	return &ChannelConfig{
		Repo: repo,
	}, nil
}

func (channel *ChannelConfig) InitConsumer(ch *amqp.Channel) {
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
			log.Printf("Received a query message with messageID %s : %s\n", data.MessageId, data.Body)
			log.Printf("ReplyTo: %s, CorrelationID: %s\n", data.ReplyTo, data.CorrelationId)
			if data.MessageId == "Post-Service" {
				type Body struct {
					Id      string      `json:"id"`
					Payload []string    `json:"payload"`
				}
				var b *Body
				err := json.Unmarshal(data.Body, &b)
				if err != nil {
					FailOnError(err, "Failed unmarshalling data body")
				}
				if b.Payload == nil {
					comments, err := channel.Repo.GetCommentsByPostId(b.Id)
					FailOnError(err, "Failed getting comments from db")

					body, err := json.Marshal(comments)

					err = ch.Publish(
						data.ReplyTo,
						"",
						false,
						false,
						amqp.Publishing{
							ContentType: 	"text/plain",
							CorrelationId: 	data.CorrelationId,
							MessageId: 		"Comment-Service",
							Body: 			body,
						})
					FailOnError(err, "Failed to publish a message")
				}
			}
		}
	}()

	go func() {
		for data := range commands {
			log.Printf("Received a command message with messageID %s : %s", data.MessageId, data.Body)
			if data.MessageId == "Post-Service" {
				var comment model.Comment
				err := json.Unmarshal(data.Body, &comment)
				if err != nil {
					log.Fatalln(err)
				}

				log.Printf("Comment: %v", comment)

				var commentEvent database.CommentEvent
				if comment.Event == "CreateComment" {
					commentEvent = database.CommentEvent{
						EventTime: time.Now().Format("2006-01-02 15:04:05"),
						EventType: comment.Event,
						PostID:    comment.PostID,
						Comment:   comment.Comment,
						CreatedBy: comment.CreatedBy,
					}
				} else {
					commentEvent = database.CommentEvent{
						EventTime: time.Now().Format("2006-01-02 15:04:05"),
						EventType: comment.Event,
						CommentID: comment.ID,
						PostID:    comment.PostID,
						Comment:   comment.Comment,
						CreatedBy: comment.CreatedBy,
					}
				}

				commentDB, err := channel.Repo.CreateComment(commentEvent)
				if err != nil {
					log.Fatalln(err)
				}

				log.Printf("CommentDB: %v", commentDB)
			}
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
