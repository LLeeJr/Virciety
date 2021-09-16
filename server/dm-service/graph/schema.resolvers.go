package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"dm-service/database"
	"dm-service/graph/generated"
	"dm-service/graph/model"
	"strings"
	"time"
)

func (r *mutationResolver) CreateDm(ctx context.Context, input *model.CreateDmRequest) (*model.Dm, error) {
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

	// push event to db
	dm, err := r.repo.CreateDm(ctx, dmEvent)
	if err != nil {
		return nil, err
	}

	// post message on queue
	r.publisher.AddMessageToEvent(dmEvent, "Dm-Service")
	r.publisher.AddMessageToCommand("Dm-Service")
	//msg := fmt.Sprintf("created new DM: %s <-> %s", dmEvent.From, dmEvent.To)
	//r.publisher.AddMessageToQuery()

	r.dmChan <- dm

	return dm, nil
}

func (r *queryResolver) GetDms(ctx context.Context) ([]*model.Dm, error) {
	dms, err := r.repo.GetDms(ctx)
	if err != nil {
		return nil, err
	}
	return dms, nil
}

func (r *queryResolver) GetChat(ctx context.Context, input model.GetChatRequest) ([]*model.Dm, error) {
	return r.repo.GetChat(ctx, input.User1, input.User2)
}

func (r *queryResolver) GetOpenChats(ctx context.Context, userName string) ([]*model.Chat, error) {
	chats, err := r.repo.GetOpenChats(ctx, userName)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *subscriptionResolver) DmAdded(ctx context.Context) (<-chan *model.Dm, error) {
	return r.dmChan, nil
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
