package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"posts-service/database"
	"posts-service/graph/generated"
	"posts-service/graph/model"
	"time"
)

func (r *mutationResolver) CreatePost(ctx context.Context, newPost *model.CreatePostRequest) (*model.Post, error) {
	postEvent := database.PostEvent{
		EventTime:   time.Now().String(),
		EventType:   "CreatePost",
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

	return post, nil
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
