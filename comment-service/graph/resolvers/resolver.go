package resolvers

import (
	"comment-service/database"
	message_queue "comment-service/message-queue"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	repo          database.Repository
	producerQueue message_queue.Publisher
}

func NewResolver(repo database.Repository, producerQueue message_queue.Publisher) *Resolver {
	return &Resolver{
		repo:          repo,
		producerQueue: producerQueue,
	}
}
