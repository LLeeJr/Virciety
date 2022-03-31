package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"notifs-service/database"
	"notifs-service/graph/generated"
	"notifs-service/graph/model"

	"github.com/google/uuid"
)

// SetReadStatus updates the read-status of a notification for a given id
func (r *mutationResolver) SetReadStatus(ctx context.Context, id string, status bool) (*model.Notif, error) {
	notification, err := r.repo.GetNotification(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = r.repo.UpdateNotification(ctx, id, status)
	if err != nil {
		return nil, err
	}

	notification.Read = status
	return notification, nil
}

// GetNotifsByReceiver returns the most-recent ten messages for a given userName
func (r *queryResolver) GetNotifsByReceiver(ctx context.Context, receiver string) ([]*model.Notif, error) {
	return r.repo.GetNotifsByReceiver(ctx, receiver)
}

// NotifAdded handles the subscription for real-time-notification-functionality. Users subscribe by providing their userName
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
