package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"event-service/database"
	"event-service/graph/generated"
	"event-service/graph/model"
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
		Members:     make([]string, 0),
		Attending:   make([]string, 0),
	}

	// save event in database
	eventModel, timeType, err := r.repo.CreateEvent(event)
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
		Members:     edit.Members,
		Description: edit.Description,
		StartDate:   edit.StartDate,
		EndDate:     edit.EndDate,
		Location:    edit.Location,
		Attending:   edit.Attending,
	}

	// save event in database
	ok, timeType, err := r.repo.EditEvent(event)
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
	ok, err := r.repo.RemoveEvent(event)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) SubscribeEvent(ctx context.Context, subscribe model.EditEventRequest) (string, error) {
	// process the data and create new event
	event := database.Event{
		EventID:     subscribe.EventID,
		EventTime:   time.Now(),
		EventType:   "SubscribeEvent",
		Title:       subscribe.Title,
		Members:     subscribe.Members,
		Description: subscribe.Description,
		StartDate:   subscribe.StartDate,
		EndDate:     subscribe.EndDate,
		Location:    subscribe.Location,
		Attending:   subscribe.Attending,
	}

	// save event in database
	ok, err := r.repo.SubscribeEvent(event)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) AttendEvent(ctx context.Context, attend model.EditEventRequest) (string, error) {
	// process the data and create new event
	event := database.Event{
		EventID:     attend.EventID,
		EventTime:   time.Now(),
		EventType:   "AttendEvent",
		Title:       attend.Title,
		Members:     attend.Members,
		Description: attend.Description,
		StartDate:   attend.StartDate,
		EndDate:     attend.EndDate,
		Location:    attend.Location,
		Attending:   attend.Attending,
	}

	// save event in database
	ok, err := r.repo.AttendEvent(event)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
