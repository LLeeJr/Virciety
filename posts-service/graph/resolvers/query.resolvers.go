package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"posts-service/graph/generated"
	"posts-service/graph/model"
)

func (r *queryResolver) GetPosts(ctx context.Context) ([]*model.Post, error) {
	return r.repo.GetPosts()
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
