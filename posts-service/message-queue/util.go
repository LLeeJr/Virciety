package message_queue

import (
	"log"
	"posts-service/database"
)

type RabbitMsg struct {
	QueueName string             `json:"queueName"`
	PostEvent database.PostEvent `json:"postEvent"`
}

type ChannelConfig struct {
	QChan chan RabbitMsg
	EChan chan RabbitMsg
	Repo  database.Repository
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
