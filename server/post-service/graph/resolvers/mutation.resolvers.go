package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"posts-service/database"
	"posts-service/graph/generated"
	"posts-service/graph/model"
	"posts-service/util"
	"time"
)

func (r *mutationResolver) CreatePost(ctx context.Context, newPost model.CreatePostRequest) (*model.Post, error) {
	created := time.Now().Format("2006-01-02 15:04:05")

	postEvent := database.PostEvent{
		EventTime:   created,
		EventType:   "CreatePost",
		PostID:      created + "__" + newPost.Username,
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

	r.postChan <- post

	// put event on queue for notifications
	// r.producerQueue.AddMessageToQuery(postEvent)

	return post, nil
}

func (r *mutationResolver) EditPost(ctx context.Context, edit model.EditPostRequest) (string, error) {
	// get post data out of current saved posts and remove it out of the list
	_, post := r.repo.GetPostById(edit.ID)
	if post == nil {
		errMsg := "no post with id " + edit.ID + " found"
		return "failed", errors.New(errMsg)
	}

	// process the data and create new post event
	post.Description = edit.NewDescription

	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "EditPost",
		PostID:      post.ID,
		Username:    post.Username,
		Description: post.Description,
		FileID:      post.Data.Name,
		LikedBy:     post.LikedBy,
		Comments:    post.Comments,
	}

	// save event in database
	ok, err := r.repo.EditPost(postEvent)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) RemovePost(ctx context.Context, removeID string) (string, error) {
	// get post data out of current saved posts
	index, post := r.repo.GetPostById(removeID)
	if post == nil {
		errMsg := "no post with id " + removeID + " found"
		return "failed", errors.New(errMsg)
	}

	// process the data and create new post event
	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "RemovePost",
		PostID:      post.ID,
		Username:    post.Username,
		Description: post.Description,
		FileID:      post.Data.Name,
		LikedBy:     make([]string, 0),
		Comments:    make([]string, 0),
	}

	// save event in database
	ok, err := r.repo.RemovePost(postEvent, index)
	if err != nil {
		return ok, err
	}

	// put event on queue for notifications to remove all notification events for this post
	// put event on queue for comments to remove all comment events for this post
	// r.producerQueue.AddMessageToQuery(postEvent)

	return ok, nil
}

func (r *mutationResolver) LikePost(ctx context.Context, like model.LikePostRequest) ([]string, error) {
	// get post data out of current saved posts and remove it out of the list
	_, post := r.repo.GetPostById(like.ID)
	if post == nil {
		errMsg := "no post with id " + like.ID + " found"
		return nil, errors.New(errMsg)
	}

	// process the data and create new post event
	event := "LikePost"
	index := util.Search(post.LikedBy, like.Username)

	// unlike the post
	if index != -1 {
		event = "UnlikePost"
		post.LikedBy = append(post.LikedBy[:index], post.LikedBy[index+1:]...)
	} else {
		post.LikedBy = append(post.LikedBy, like.Username)
	}

	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   event,
		PostID:      post.ID,
		Username:    post.Username,
		Description: post.Description,
		FileID:      post.Data.Name,
		LikedBy:     post.LikedBy,
		Comments:    post.Comments,
	}

	// save event in database
	err := r.repo.LikePost(postEvent)
	if err != nil {
		return nil, err
	}

	// put event on queue for notifications
	// r.producerQueue.AddMessageToQuery(postEvent)

	return post.LikedBy, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
