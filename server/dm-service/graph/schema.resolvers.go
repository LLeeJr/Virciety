package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"dm-service/database"
	"dm-service/graph/generated"
	"dm-service/graph/model"
	"github.com/google/uuid"
	"time"
)

func (r *mutationResolver) CreateDm(ctx context.Context, msg string, username string, roomName string) (*model.Dm, error) {

	r.mu.Lock()
	room := r.Rooms[roomName]
	if room == nil {
		room = &Chatroom{
			Name: roomName,
			Observers: map[string]struct {
				Username string
				Message  chan *model.Dm
			}{},
		}
		r.Rooms[roomName] = room
	}
	r.mu.Unlock()
	
	id := uuid.NewString()
	dmEvent := database.DmEvent{
		EventTime: time.Now().Format("2006-01-02 15:04:05"),
		EventType: "CreateDm",
		DmID:      id,
		From:      username,
		Msg:       msg,
	}

	// push event to db
	dm, err := r.repo.CreateDm(ctx, dmEvent)
	if err != nil {
		return nil, err
	}

	room.Messages = append(room.Messages, dm)
	r.mu.Lock()
	for _, observer := range room.Observers {
		observer.Message <- dm
	}
	r.mu.Unlock()
	
	// post message on queue
	r.publisher.AddMessageToEvent(dmEvent, "Dm-Service")
	r.publisher.AddMessageToCommand("Dm-Service")
	//msg := fmt.Sprintf("created new DM: %s <-> %s", dmEvent.From, dmEvent.To)
	//r.publisher.AddMessageToQuery()

	return dm, nil
}

func (r *queryResolver) GetRoom(ctx context.Context, name string) (*model.Chatroom, error) {
	r.mu.Lock()
	room := r.Rooms[name]
	if room == nil {
		room = &Chatroom{
			Name: name,
			Observers: map[string]struct {
				Username string
				Message  chan *model.Dm
			}{},
		}
		r.Rooms[name] = room
	}
	r.mu.Unlock()

	chatroom := &model.Chatroom{
		Name:     room.Name,
		Messages: room.Messages,
	}

	return chatroom, nil
}

// GetDms Deprecated
func (r *queryResolver) GetDms(ctx context.Context) ([]*model.Dm, error) {
	dms, err := r.repo.GetDms(ctx)
	if err != nil {
		return nil, err
	}
	return dms, nil
}

// GetChat Deprecated
func (r *queryResolver) GetChat(ctx context.Context, input model.GetChatRequest) ([]*model.Dm, error) {
	return r.repo.GetChat(ctx, input.User1, input.User2)
}

// GetOpenChats Deprecated
func (r *queryResolver) GetOpenChats(ctx context.Context, userName string) ([]*model.Chat, error) {
	chats, err := r.repo.GetOpenChats(ctx, userName)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *subscriptionResolver) DmAdded(ctx context.Context, roomName string) (<-chan *model.Dm, error) {
	r.mu.Lock()
	room := r.Rooms[roomName]
	if room == nil {
		room = &Chatroom{
			Name: roomName,
			Observers: map[string]struct {
				Username string
				Message  chan *model.Dm
			}{},
		}
		r.Rooms[roomName] = room
	}
	r.mu.Unlock()
	
	id := uuid.NewString()
	events := make(chan *model.Dm, 1)
	
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(room.Observers, id)
		r.mu.Unlock()
	}()

	r.mu.Lock()
	room.Observers[id] = struct {
		Username string
		Message  chan *model.Dm
	}{Username: getUsername(ctx), Message: events}
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

type Chatroom struct {
	Name      string
	Messages  []*model.Dm
	Observers map[string]struct {
		Username string
		Message  chan *model.Dm
	}
}