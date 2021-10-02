package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"posts-service/database"
	"posts-service/graph/generated"
	"posts-service/graph/model"
	"time"
)

func (r *mutationResolver) CreatePost(ctx context.Context, newPost model.CreatePostRequest) (*model.Post, error) {
	created := time.Now().Format("2006-01-02 15:04:05")

	postEvent := database.PostEvent{
		EventTime:   created,
		EventType:   "CreatePost",
		Username:    newPost.Username,
		Description: newPost.Description,
		FileID:      "",
		LikedBy:     make([]string, 0),
		Comments:    make([]string, 0),
	}

	// save event in database
	post, err := r.repo.CreatePost(postEvent, newPost.Data)
	if err != nil {
		return nil, err
	}

	// r.postChan <- post

	// put event on queue for notifications
	// r.producerQueue.AddMessageToQuery(postEvent)

	return post, nil
}

func (r *mutationResolver) EditPost(ctx context.Context, edit model.EditPostRequest) (string, error) {
	// process the data and create new post event
	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "EditPost",
		PostID:      edit.ID,
		Description: edit.NewDescription,
		LikedBy:     edit.LikedBy,
		Comments:    edit.Comments,
	}

	// save event in database
	ok, err := r.repo.EditPost(postEvent)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) RemovePost(ctx context.Context, remove model.RemovePostRequest) (string, error) {
	// process the data and create new post event
	postEvent := database.PostEvent{
		EventTime: time.Now().Format("2006-01-02 15:04:05"),
		EventType: "RemovePost",
		PostID:    remove.ID,
		FileID:    remove.FileID,
	}

	// save event in database
	ok, err := r.repo.RemovePost(postEvent)
	if err != nil {
		return ok, err
	}

	// put event on queue for notifications to remove all notification events for this post
	// put event on queue for comments to remove all comment events for this post
	// r.producerQueue.AddMessageToQuery(postEvent)

	return ok, nil
}

func (r *mutationResolver) LikePost(ctx context.Context, like model.LikePostRequest) (string, error) {
	// process the data and create new post event
	event := "LikePost"

	if !like.Liked {
		event = "UnlikePost"
	}

	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   event,
		PostID:      like.ID,
		Description: like.Description,
		LikedBy:     like.NewLikedBy,
		Comments:    like.Comments,
	}

	// save event in database
	err := r.repo.LikePost(postEvent)
	if err != nil {
		return "failed", err
	}

	// put event on queue for notifications
	// r.producerQueue.AddMessageToQuery(postEvent)

	return "success", nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
