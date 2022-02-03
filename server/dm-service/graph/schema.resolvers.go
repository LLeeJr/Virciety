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

func (r *mutationResolver) CreateDm(ctx context.Context, msg string, userName string, roomName string, roomID string) (*model.Dm, error) {
	r.mu.Lock()
	room := r.Rooms[roomName]

	// if no room was found in current server session, try retrieving from db
	if room == nil {
		chatroom, err := r.repo.GetRoom(ctx, roomName, roomID)
		if err != nil {
			// room does not exist in db
			r.mu.Unlock()
			return nil, err
		}

		if chatroom != nil {
			room = &Chatroom{
				Name:   chatroom.Name,
				Member: chatroom.Member,
				Id:     chatroom.ID,
				Owner:  chatroom.Owner,
				Observers: map[string]struct {
					Username string
					Message  chan *model.Dm
				}{},
			}
			r.Rooms[roomName] = room
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

func (r *mutationResolver) CreateRoom(ctx context.Context, input model.CreateRoom) (*model.Chatroom, error) {
	r.mu.Lock()
	room := r.Rooms[input.Name]

	// if no room was found in current server session, try retrieving from db
	if r.Rooms[input.Name] == nil {
		chatroom, err := r.repo.GetRoom(ctx, input.Name, "")
		if err != nil {
			// no room was found in db, so create one
			roomEvent := database.ChatroomEvent{
				EventType:  "CreateRoom",
				Member:     input.Member,
				MemberSize: len(input.Member),
				Name:       input.Name,
				Owner:      input.Owner,
			}
			createRoom, createRoomErr := r.repo.CreateRoom(ctx, roomEvent)
			if createRoomErr != nil {
				r.mu.Unlock()
				return nil, createRoomErr
			}
			room = &Chatroom{
				Id:     createRoom.ID,
				Name:   createRoom.Name,
				Member: createRoom.Member,
				Owner:  createRoom.Owner,
				Observers: map[string]struct {
					Username string
					Message  chan *model.Dm
				}{},
			}
			r.Rooms[createRoom.Name] = room
		} else {
			// no error occurred, so a room was found in db
			room = &Chatroom{
				Id:     chatroom.ID,
				Name:   chatroom.Name,
				Member: chatroom.Member,
				Owner:  chatroom.Owner,
				Observers: map[string]struct {
					Username string
					Message  chan *model.Dm
				}{},
			}
			r.Rooms[chatroom.Name] = room
		}
	}

	r.mu.Unlock()
	return &model.Chatroom{
		ID:     room.Id,
		Member: room.Member,
		Name:   room.Name,
		Owner:  room.Owner,
	}, nil
}

func (r *mutationResolver) DeleteRoom(ctx context.Context, remove model.RemoveRoom) (string, error) {
	r.mu.Lock()
	room := r.Rooms[remove.RoomName]

	// if room does not exist in cache look for the room in db
	if r.Rooms[remove.RoomName] == nil {
		r.mu.Unlock()
		chatroom, err := r.repo.GetRoom(ctx, remove.RoomName, remove.ID)
		if err != nil {
			// room does not exist in db
			return "", err
		}
		if chatroom.Owner != remove.UserName {
			err = errors.New("given user is no owner of the requested room")
			return "", err
		}
		msg, err := r.repo.DeleteRoom(ctx, chatroom.ID)
		if err != nil {
			return "", err
		}
		return msg, err
	}

	// room does exist in server cache
	if room.Owner != remove.UserName {
		r.mu.Unlock()
		err := errors.New("given user is no owner of the requested room")
		return "", err
	}

	msg, err := r.repo.DeleteRoom(ctx, room.Id)
	if err != nil {
		r.mu.Unlock()
		return "", err
	}

	delete(r.Rooms, room.Name)
	r.mu.Unlock()

	return msg, err
}

func (r *queryResolver) GetDirectRoom(ctx context.Context, user1 string, user2 string) (*model.Chatroom, error) {
	room, err := r.repo.GetDirectRoom(ctx, user1, user2)
	if err != nil {
		return nil, err
	}

	if r.Rooms[room.Name] == nil {
		chatroom := &Chatroom{
			Id:     room.ID,
			Name:   room.Name,
			Member: room.Member,
			Owner:  room.Owner,
			Observers: map[string]struct {
				Username string
				Message  chan *model.Dm
			}{},
		}
		r.Rooms[room.Name] = chatroom
	}

	return room, nil
}

func (r *queryResolver) GetRoom(ctx context.Context, roomName string, roomID string) (*model.Chatroom, error) {
	r.mu.Lock()
	room := r.Rooms[roomName]

	// if room does not exist in cache look for the room in db
	if r.Rooms[roomName] == nil {
		chatroom, err := r.repo.GetRoom(ctx, roomName, roomID)
		if err != nil {
			// room does not exist in db
			r.mu.Unlock()
			return nil, err
		}

		if chatroom != nil {
			// room exists in db
			room = &Chatroom{
				Id:     chatroom.ID,
				Name:   chatroom.Name,
				Member: chatroom.Member,
				Owner:  chatroom.Owner,
				Observers: map[string]struct {
					Username string
					Message  chan *model.Dm
				}{},
			}
			r.Rooms[roomName] = room
		}
	}
	r.mu.Unlock()

	return &model.Chatroom{
		ID:     room.Id,
		Member: room.Member,
		Name:   room.Name,
		Owner:  room.Owner,
	}, nil
}

func (r *queryResolver) GetRoomsByUser(ctx context.Context, userName string) ([]*model.Chatroom, error) {
	rooms, err := r.repo.GetRoomsByUser(ctx, userName)
	if err != nil {
		return nil, err
	}
	r.mu.Lock()
	ids := make([]string, 0)
	for _, chatroom := range r.Rooms {
		ids = append(ids, chatroom.Id)
	}
	for _, room := range rooms {
		if !util.Contains(ids, room.ID) {
			r.Rooms[room.Name] = &Chatroom{
				Id:     room.ID,
				Name:   room.Name,
				Member: room.Member,
				Owner:  room.Owner,
				Observers: map[string]struct {
					Username string
					Message  chan *model.Dm
				}{},
			}
		}
	}
	r.mu.Unlock()
	return rooms, nil
}

func (r *queryResolver) GetMessagesFromRoom(ctx context.Context, roomID string) ([]*model.Dm, error) {
	dms, err := r.repo.GetMessagesFromRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return dms, nil
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
	Owner     string
	Observers map[string]struct {
		Username string
		Message  chan *model.Dm
	}
}
