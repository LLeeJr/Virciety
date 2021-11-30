package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"event-service/graph/generated"
	"event-service/graph/model"
)

func (r *queryResolver) GetEvents(ctx context.Context) ([]*model.Event, error) {
	return r.repo.GetEvents()
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
