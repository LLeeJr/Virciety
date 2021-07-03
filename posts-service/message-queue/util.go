package message_queue

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
	"posts-service/database"
	"strings"
)

type RabbitMsg struct {
	QueueName string             `json:"queueName"`
	PostEvent database.PostEvent `json:"postEvent"`
	MessageId string             `json:"messageId"`
}

type ChannelConfig struct {
	QueryChan   chan RabbitMsg
	CommandChan chan RabbitMsg
	EventChan   chan RabbitMsg
	Repo        database.Repository
}

type Comment struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	LikedBy     []string `json:"likedBy"`
}

func initExchange(queueName string, ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		queueName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func initQueue(queueName string, exchangeName string, ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		exchangeName,
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
}

func convertCommentEventToPostEvent(repo database.Repository, data []byte) (*database.PostEvent, error) {
	var commentEvent database.CommentEvent

	err := json.Unmarshal(data, &commentEvent)
	if err != nil {
		return nil, err
	}

	_, post := repo.GetPostById(commentEvent.PostID)
	if post == nil {
		errMsg := "no post with id " + commentEvent.PostID + " found"
		return nil, errors.New(errMsg)
	}

	b, err := json.Marshal(&Comment{
		ID:          commentEvent.CommentID,
		Description: commentEvent.Description,
		LikedBy:     commentEvent.LikedBy,
	})
	if err != nil {
		return nil, err
	}

	post.Comments = append(post.Comments, string(b))

	info := strings.Split(post.ID, "__")

	return &database.PostEvent{
		EventTime:   commentEvent.EventTime,
		EventType:   commentEvent.EventType,
		PostID:      commentEvent.PostID,
		Username:    info[1],
		Description: post.Description,
		Data:        post.Data,
		LikedBy:     post.LikedBy,
		Comments:    post.Comments,
	}, nil
}
