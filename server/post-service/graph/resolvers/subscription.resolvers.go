package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"posts-service/graph/generated"
	"posts-service/graph/model"

	"github.com/google/uuid"
)

func (r *subscriptionResolver) NewPostCreated(ctx context.Context) (<-chan *model.Post, error) {
	fmt.Println("Subscribed")

	id := uuid.NewString()
	events := make(chan *model.Post, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		fmt.Printf("Deleted id %v\n", id)
		delete(r.observers, id)
		r.mu.Unlock()
	}()

	r.mu.Lock()
	r.observers[id] = events
	r.mu.Unlock()

	return events, nil
}

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
