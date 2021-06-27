package gateway

import "context"

type Service interface {
	CreateDMEvent(ctx context.Context, id string, msg string) (string, error)
	GetDMEvent(ctx context.Context, id string) (string, error)
}