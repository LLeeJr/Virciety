package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"posts-service/graph/generated"
	"posts-service/graph/model"
)

func (r *queryResolver) GetPosts(ctx context.Context) ([]*model.Post, error) {
	currentPosts, err := r.repo.GetPosts()
	if err != nil {
		return nil, err
	}

	return currentPosts, nil
}

func (r *queryResolver) GetData(ctx context.Context, id string) (string, error) {
	currentPosts := r.repo.GetCurrentPosts()

	for _, post := range currentPosts {
		if post.ID == id {
			return post.Data.Content, nil
		}
	}

	return "", errors.New("post with id " + id + " not found")
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
