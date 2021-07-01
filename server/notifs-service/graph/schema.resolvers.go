package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"notifs-service/database"
	"notifs-service/graph/generated"
	"notifs-service/graph/model"
	"time"
)

func (r *mutationResolver) CreateNotif(ctx context.Context, input model.CreateNotifRequest) (*model.Notif, error) {
	id := fmt.Sprintf("%s__%s__%s", input.Event, time.Now().Format("2006-01-02 15:04:05"), input.Receiver)

	notifEvent := database.NotifEvent{
		EventTime: time.Now().Format("2006-01-02 15:04:05"),
		EventType: input.Event,
		NotifId:   id,
		Receiver:  input.Receiver,
		Text:      input.Text,
	}
	log.Println("notifEvent", notifEvent)

	notif, err := r.repo.CreateNotif(ctx, notifEvent)
	if err != nil {
		return nil, err
	}
	r.notifs = append(r.notifs, notif)
	r.notifsChan <- notif

	return notif, nil
}

func (r *queryResolver) GetNotifsByReceiver(ctx context.Context, receiver string) ([]*model.Notif, error) {
	return r.repo.GetNotifsByReceiver(ctx, receiver)
}

func (r *subscriptionResolver) NotifAdded(ctx context.Context) (<-chan *model.Notif, error) {
	return r.notifsChan, nil
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
