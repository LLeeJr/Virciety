package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"event-service/database"
	"event-service/graph/generated"
	"event-service/graph/model"
	message_queue "event-service/message-queue"
	"fmt"
	"time"
)

func (r *mutationResolver) CreateEvent(ctx context.Context, newEvent model.CreateEventRequest) (*model.CreateEventResponse, error) {
	created := time.Now()

	event := database.Event{
		EventID:     "",
		EventTime:   created,
		EventType:   "CreateEvent",
		Title:       newEvent.Title,
		Host:        newEvent.Host,
		Description: newEvent.Description,
		StartDate:   newEvent.StartDate,
		EndDate:     newEvent.EndDate,
		Location:    newEvent.Location,
		Subscribers: make([]string, 0),
		Attendees:   make([]string, 0),
	}

	// save event in database
	eventModel, timeType, err := r.repo.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return &model.CreateEventResponse{
		Event: eventModel,
		Type:  timeType,
	}, nil
}

func (r *mutationResolver) EditEvent(ctx context.Context, edit model.EditEventRequest) (*model.EditEventResponse, error) {
	// process the data and create new event
	event := database.Event{
		EventID:     edit.EventID,
		EventTime:   time.Now(),
		EventType:   "EditEvent",
		Title:       edit.Title,
		Subscribers: edit.Subscribers,
		Description: edit.Description,
		StartDate:   edit.StartDate,
		EndDate:     edit.EndDate,
		Location:    edit.Location,
		Attendees:   edit.Attendees,
	}

	message := fmt.Sprintf("Changes on event %s", edit.Title)
	eventNotification := message_queue.EventNotification{
		EventId:  edit.EventID,
		Message:  message,
		EditFlag: true,
	}
	for _, attendee := range edit.Subscribers {
		eventNotification.Username = attendee
		r.producerQueue.AddMessageToEvent(eventNotification)
	}

	// save event in database
	ok, timeType, err := r.repo.EditEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return &model.EditEventResponse{
		Ok:   ok,
		Type: timeType,
	}, nil
}

func (r *mutationResolver) RemoveEvent(ctx context.Context, remove string) (string, error) {
	// process the data and create new event
	event := database.Event{
		EventID:   remove,
		EventTime: time.Now(),
		EventType: "RemoveEvent",
	}

	// save event in database
	ok, err := r.repo.RemoveEvent(ctx, event)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) SubscribeEvent(ctx context.Context, eventID string, username *string, subscribed bool) ([]string, error) {
	r.mu.Lock()
	updatedEvent, err := r.repo.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if subscribed {
		updatedEvent.Subscribers = append(updatedEvent.Subscribers, *username)
	} else {
		index := -1
		for i, subscriber := range updatedEvent.Subscribers {
			if subscriber == *username {
				index = i
				break
			}
		}

		if index == -1 {
			return nil, errors.New("username not found in subscriber list")
		}

		updatedEvent.Subscribers = append(updatedEvent.Subscribers[:index], updatedEvent.Subscribers[index+1:]...)
	}

	// process the data and create new event
	event := database.Event{
		EventID:     updatedEvent.ID,
		EventTime:   time.Now(),
		EventType:   "SubscribeEvent",
		Title:       updatedEvent.Title,
		Subscribers: updatedEvent.Subscribers,
		Description: updatedEvent.Description,
		StartDate:   updatedEvent.StartDate,
		EndDate:     updatedEvent.EndDate,
		Location:    updatedEvent.Location,
		Attendees:   updatedEvent.Attendees,
	}

	// save event in database
	_, err = r.repo.InsertEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	r.mu.Unlock()

	return updatedEvent.Subscribers, nil
}

func (r *mutationResolver) AttendEvent(ctx context.Context, eventID string, username string, left bool) ([]string, error) {
	r.mu.Lock()
	// save event in database
	attendees, err := r.repo.AttendEvent(ctx, eventID, username, left)
	if err != nil {
		return nil, err
	}
	r.mu.Unlock()

	return attendees, nil
}

func (r *mutationResolver) AddUserData(ctx context.Context, userData model.UserDataRequest) (*model.UserData, error) {
	return r.repo.InsertUserData(ctx, userData)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
