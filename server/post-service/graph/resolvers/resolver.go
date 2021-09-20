package resolvers

import (
	"posts-service/database"
	"posts-service/graph/model"
	messagequeue "posts-service/message-queue"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	repo          database.Repository
	producerQueue messagequeue.Publisher
	postChan      chan *model.Post
}

func NewResolver(repo database.Repository, producerQueue messagequeue.Publisher) *Resolver {
	return &Resolver{
		repo:          repo,
		producerQueue: producerQueue,
		postChan:      make(chan *model.Post),
	}
}
