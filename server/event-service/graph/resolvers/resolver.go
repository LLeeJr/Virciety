package resolvers

import (
	"event-service/database"
	messagequeue "event-service/message-queue"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	repo          database.Repository
	producerQueue messagequeue.Publisher
	mu            sync.Mutex
}

func NewResolver(repo database.Repository, producerQueue messagequeue.Publisher) *Resolver {
	return &Resolver{
		repo:          repo,
		producerQueue: producerQueue,
	}
}
