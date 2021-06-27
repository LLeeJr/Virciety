package message

import "context"

type Service interface {
	CreateMessage(ctx context.Context, id string, msg string) (string, error)
	GetMessage(ctx context.Context, id string) (string, error)
}