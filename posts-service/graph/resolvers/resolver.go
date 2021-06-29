package resolvers

import (
	"posts-service/database"
	"posts-service/graph/model"
	message_queue "posts-service/message-queue"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	repo         database.Repository
	currentPosts []*model.Post
	messageQueue message_queue.MessageQueue
}

func NewResolver(repo database.Repository, messageQueue message_queue.MessageQueue) *Resolver {
	return &Resolver{
		repo:         repo,
		messageQueue: messageQueue,
		currentPosts: make([]*model.Post, 0),
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
