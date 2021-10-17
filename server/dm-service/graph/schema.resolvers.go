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
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *mutationResolver) CreateDm(ctx context.Context, msg string, userName string, roomName string) (*model.Dm, error) {
	r.mu.Lock()
	room := r.Rooms[roomName]

	// if room does not exist yet, create a new one
	if room == nil {
		chatroom, err := r.repo.GetRoom(ctx, roomName)
		if err != nil {
			return nil, err
		}

		if chatroom != nil {
			room = &Chatroom{
				Name:   chatroom.Name,
				Member: chatroom.Member,
				Id:     chatroom.ID,
				Observers: map[string]struct {
					Username string
					Message  chan *model.Dm
				}{},
			}
		} else {
			// room exists neither in db nor in server cache
			room = &Chatroom{
				Name:   roomName,
				Member: []string{userName},
				Observers: map[string]struct {
					Username string
					Message  chan *model.Dm
				}{},
			}
			roomEvent := database.ChatroomEvent{
				EventType: "CreateRoom",
				Member:    room.Member,
				Name:      room.Name,
			}

			createdRoom, err := r.repo.CreateRoom(ctx, roomEvent)
			if err != nil {
				return nil, err
			}
			// update repos room map since the current room now has an id
			room.Id = createdRoom.ID
		}
		r.Rooms[roomName] = room
	}

	// if the user is no member of the chatroom yet, add the user as member
	if !util.Contains(room.Member, userName) {
		room.Member = append(room.Member, userName)
		_, err := r.repo.UpdateRoom(ctx, &model.Chatroom{
			ID:     room.Id,
			Member: room.Member,
			Name:   room.Name,
		})
		if err != nil {
			return nil, err
		}
	}
	r.mu.Unlock()

	dmEvent := database.DmEvent{
		ChatroomId: room.Id,
		CreatedAt:  time.Now(),
		CreatedBy:  userName,
		EventType:  "CreateDm",
		Msg:        msg,
	}

	// push event to db
	dm, err := r.repo.CreateDm(ctx, dmEvent)
	if err != nil {
		return nil, err
	}

	// write new dm to observers of the respective chatroom
	room.Messages = append(room.Messages, dm)
	r.mu.Lock()
	for _, observer := range room.Observers {
		observer.Message <- dm
	}
	r.mu.Unlock()

	// post message on queue
	r.publisher.AddMessageToEvent(dmEvent, "Dm-Service")
	r.publisher.AddMessageToCommand("Dm-Service")
	//r.publisher.AddMessageToQuery()

	return dm, nil
}

func (r *queryResolver) GetRoom(ctx context.Context, name string) (*model.Chatroom, error) {

	r.mu.Lock()
	room := r.Rooms[name]

	// if room does not exist in cache look for the room in db
	if r.Rooms[name] == nil {
		chatroom, err := r.repo.GetRoom(ctx, name)
		if err != nil {
			return nil, err
		}

		if chatroom != nil {
			// room exists in db
			room = &Chatroom{
				Id:        chatroom.ID,
				Name:      chatroom.Name,
				Member:    chatroom.Member,
				Observers: map[string]struct {
					Username string
					Message  chan *model.Dm
				}{},
			}
			r.Rooms[name] = room
		} else {
			errorMsg := fmt.Sprint("no existing room with given name: ", name)
			return nil, errors.New(errorMsg)
		}
	}
	r.mu.Unlock()

	return &model.Chatroom{
		ID:        room.Id,
		Name:      room.Name,
		Member:    room.Member,
	}, nil
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
				Member: chatroom.Member,
				Name:   chatroom.Name,
			})
		}
	}
	r.mu.Unlock()

	if len(rooms) == 0 {
		return nil, errors.New("no rooms available for this user")
	}

	return rooms, nil
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
	Id        string
	Name      string
	Messages  []*model.Dm
	Member    []string
	Observers map[string]struct {
		Username string
		Message  chan *model.Dm
	}
}
