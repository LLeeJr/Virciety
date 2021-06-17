package message

import "context"

type Service interface {
	CreateMessage(ctx context.Context, msg string, from string, to string) (string, error)
	GetMessage(ctx context.Context, id string) (string, error)
}