package graph

import (
	"context"
	"dm-service/database"
	"dm-service/queue"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ckey string

// NewResolver creates a new Resolver struct with a given rabbitMQ-publisher and an instance of the repository
func NewResolver(repo database.Repository, publisher queue.Publisher) *Resolver {
	return &Resolver{
		Rooms: map[string]*Chatroom{},
		repo:   repo,
		publisher: publisher,
	}
}

type Resolver struct {
	Rooms map[string]*Chatroom
	mu sync.Mutex
	repo      database.Repository
	publisher queue.Publisher
}

func getUsername(ctx context.Context) string {
	if username, ok := ctx.Value(ckey("username")).(string); ok {
		return username
	}
	return ""
}