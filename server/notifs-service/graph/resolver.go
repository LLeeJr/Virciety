package graph

import (
	"notifs-service/database"
	"notifs-service/queue"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func NewResolver(repo database.Repository, publisher queue.Publisher) *Resolver {
	return &Resolver{
		repo:       repo,
		publisher:  publisher,
	}
}

type Resolver struct {
	repo		database.Repository
	publisher	queue.Publisher
	mu sync.Mutex
}
