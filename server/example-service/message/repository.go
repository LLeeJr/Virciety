package message

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
)

var RepoErr = errors.New("unable to handle Repository Request")
var messages []Message

type repo struct {
	db 	   *sql.DB
	logger log.Logger
}

func (r *repo) CreateMessage(ctx context.Context, msg Message) error {
	if msg.From == "" || msg.To == "" {
		return RepoErr
	}

	fmt.Printf("Message from %s to %s: %s", msg.From, msg.To, msg.Msg)
	messages = append(messages, msg)
	return nil
}

func (r *repo) GetMessage(ctx context.Context, id string) (string, error) {
	for _, message := range messages {
		if message.ID == id {
			fmt.Printf("Message with ID %s: %s", id, message.Msg)
			return message.Msg, nil
		}
	}
	return "Message not found", nil
}

func NewRepo(db *sql.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repository", "sql"),
	}
}
