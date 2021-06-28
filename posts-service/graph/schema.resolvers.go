package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"posts-service/database"
	"posts-service/graph/generated"
	"posts-service/graph/model"
	"strings"
	"time"
)

func (r *mutationResolver) CreatePost(ctx context.Context, newPost *model.CreatePostRequest) (*model.Post, error) {
	created := time.Now().String()

	postEvent := database.PostEvent{
		EventTime:   created,
		EventType:   "CreatePost",
		PostID:      created + "__" + newPost.Username,
		Username:    newPost.Username,
		Description: newPost.Description,
		Data:        newPost.Data,
		LikedBy:     make([]string, 0),
		Comments:    make([]string, 0),
	}

	// put event on queue

	// save event in database
	post, err := r.repo.CreatePost(postEvent)
	if err != nil {
		return nil, err
	}

	// add to currentPosts in resolver
	r.currentPosts = append(r.currentPosts, *post)

	return post, nil
}

func (r *mutationResolver) EditPost(ctx context.Context, edit *model.EditPostRequest) (string, error) {
	// Get post data out of current saved posts and remove it out of the list
	index, post := GetPostByID(r.currentPosts, edit.ID)
	if post == nil {
		errMsg := "no post with id " + edit.ID + " found"
		return "failed", errors.New(errMsg)
	}

	r.currentPosts = append(r.currentPosts[:index], r.currentPosts[index+1:]...)

	// Process the data and create new post event
	info := strings.Split(post.ID, "__")

	post.Description = edit.NewDescription

	postEvent := database.PostEvent{
		EventTime:   time.Now().String(),
		EventType:   "EditPost",
		PostID:      post.ID,
		Username:    info[1],
		Description: post.Description,
		Data:        post.Data,
		LikedBy:     post.LikedBy,
		Comments:    post.Comments,
	}

	// Save the new event in database
	ok, err := r.repo.EditPost(postEvent)
	if err != nil {
		return ok, err
	}

	// Add post to currentPosts
	r.currentPosts = append(r.currentPosts, *post)

	return ok, nil
}

func (r *queryResolver) GetPosts(ctx context.Context) ([]*model.Post, error) {
	return r.repo.GetPosts()
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
