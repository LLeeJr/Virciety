package message

import "context"

type Message struct {
	ID string `json:"id,omitempty"`
	Msg string `json:"msg"`
}

type Repository interface {
	CreateMessage(ctx context.Context, msg Message) error
	GetMessage(ctx context.Context, id string) (string, error)
}