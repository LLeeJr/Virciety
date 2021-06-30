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
	currentPosts  []*model.Post
	producerQueue messagequeue.ProducerQueue
}

func NewResolver(repo database.Repository, producerQueue messagequeue.ProducerQueue) *Resolver {
	return &Resolver{
		repo:          repo,
		producerQueue: producerQueue,
		currentPosts:  make([]*model.Post, 0),
	}
}

func GetPostByID(currentPosts []*model.Post, id string) (int, *model.Post) {
	for i, post := range currentPosts {
		if post.ID == id {
			return i, post
		}
	}

	return -1, nil
}
