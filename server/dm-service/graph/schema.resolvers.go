package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"dm-service/database"
	"dm-service/graph/generated"
	"dm-service/graph/model"
	"log"
	"strings"
	"time"
)

func (r *mutationResolver) CreateDm(ctx context.Context, input *model.CreateDmRequest) (*model.Dm, error) {
	log.Println("input", input)

	data := strings.Split(input.ID, "__")
	dmEvent := database.DmEvent{
		EventTime: time.Now().Format("2006-01-02 15:04:05"),
		EventType: "CreateDm",
		DmID:      input.ID,
		From:      data[0],
		To:        data[2],
		Time:      data[1],
		Msg:       input.Msg,
	}
	log.Println("dmEvent", dmEvent)

	dm, err := r.repo.CreateDm(ctx, dmEvent)
	if err != nil {
		return nil, err
	}
	r.dms = append(r.dms, dm)

	r.dmsChan <- dm

	return dm, nil
}

func (r *queryResolver) GetDms(ctx context.Context) ([]*model.Dm, error) {
	return r.repo.GetDms(ctx)
}

func (r *queryResolver) GetDmsByFromTo(ctx context.Context, input model.GetByFromToRequest) ([]*model.Dm, error) {
	return r.repo.GetDmsByFromTo(ctx, input.From, input.To)
}

func (r *subscriptionResolver) DmAdded(ctx context.Context) (<-chan *model.Dm, error) {
	return r.dmsChan, nil
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
