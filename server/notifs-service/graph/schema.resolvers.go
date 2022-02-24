package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/google/uuid"
	"notifs-service/database"
	"notifs-service/graph/generated"
	"notifs-service/graph/model"
)

func (r *queryResolver) GetNotifsByReceiver(ctx context.Context, receiver string) ([]*model.Notif, error) {
	return r.repo.GetNotifsByReceiver(ctx, receiver)
}

func (r *subscriptionResolver) NotifAdded(ctx context.Context, userName string) (<-chan *model.Notif, error) {
	r.mu.Lock()
	subscription := r.repo.GetSubscriptions()[userName]
	if subscription == nil {
		subscription = &database.Message{
			Observers: map[string]struct {
				Message chan *model.Notif
			}{},
		}
		r.repo.AddSubscription(userName, subscription)
	}
	r.mu.Unlock()

	id := uuid.NewString()
	events := make(chan *model.Notif, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(subscription.Observers, id)
		r.mu.Unlock()
	}()

	r.mu.Lock()
	subscription.Observers[id] = struct{ Message chan *model.Notif }{Message: events}
	r.mu.Unlock()

	return events, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type queryResolver struct{ *Resolver }

type subscriptionResolver struct {
	*Resolver
}
