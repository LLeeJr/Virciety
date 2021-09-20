package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"posts-service/graph/generated"
	"posts-service/graph/model"
)

func (r *subscriptionResolver) NewPostCreated(ctx context.Context) (<-chan *model.Post, error) {
	return r.postChan, nil
}

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
