package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"dm-service/database"
	"dm-service/graph/generated"
	"dm-service/graph/model"
	"dm-service/util"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (r *mutationResolver) CreateDm(ctx context.Context, msg string, userName string, roomName string) (*model.Dm, error) {
	r.mu.Lock()
	room := r.Rooms[roomName]
	if room == nil {
		room = &Chatroom{
			Name:   roomName,
			Member: []string{userName},
			Observers: map[string]struct {
				Username string
				Message  chan *model.Dm
			}{},
		}
		r.Rooms[roomName] = room
	}
	if !util.Contains(room.Member, userName) {
		room.Member = append(room.Member, userName)
	}
	r.mu.Unlock()

	id := uuid.NewString()
	dmEvent := database.DmEvent{
		EventTime: time.Now(),
		EventType: "CreateDm",
		DmID:      id,
		From:      userName,
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
			Name:   name,
			Member: make([]string, 0),
			Observers: map[string]struct {
				Username string
				Message  chan *model.Dm
			}{},
		}
		r.Rooms[name] = room
	}
	r.mu.Unlock()

	chatroom := &model.Chatroom{
		Member:   room.Member,
		Messages: room.Messages,
		Name:     room.Name,
	}

	return chatroom, nil
}

func (r *queryResolver) GetRoomsByUser(ctx context.Context, userName string) ([]*model.Chatroom, error) {
	r.mu.Lock()
	rooms := make([]*model.Chatroom, 0)
	if len(r.Rooms) == 0 {
		return nil, errors.New("no rooms available")
	}
	for _, chatroom := range r.Rooms {
		if util.Contains(chatroom.Member, userName) {
			rooms = append(rooms, &model.Chatroom{
				Member:   chatroom.Member,
				Messages: chatroom.Messages,
				Name:     chatroom.Name,
			})
		}
	}
	r.mu.Unlock()

	if len(rooms) == 0 {
		return nil, errors.New("no rooms available for this user")
	}

	return rooms, nil
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
	Member    []string
	Observers map[string]struct {
		Username string
		Message  chan *model.Dm
	}
}
