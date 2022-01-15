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

func (r *mutationResolver) CreateEvent(ctx context.Context, newEvent model.CreateEventRequest) (*model.Event, error) {
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
	}

	// save event in database
	eventModel, err := r.repo.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	return eventModel, nil
}

func (r *mutationResolver) EditEvent(ctx context.Context, edit model.EditEventRequest) (string, error) {
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
	}

	// save event in database
	ok, err := r.repo.EditEvent(event)
	if err != nil {
		return ok, err
	}

	return ok, nil
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

func (r *mutationResolver) JoinEvent(ctx context.Context, joinEvent model.JoinEventRequest) (string, error) {
	// process the data and create new event
	event := database.Event{
		EventID:     joinEvent.EventID,
		EventTime:   time.Now(),
		EventType:   "AddedMember",
		Title:       joinEvent.Title,
		Members:     joinEvent.NewMembers,
		Description: joinEvent.Description,
		StartDate:   joinEvent.StartDate,
		EndDate:     joinEvent.EndDate,
		Location:    joinEvent.Location,
	}

	// save event in database
	ok, err := r.repo.AddedMember(event)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) LeaveEvent(ctx context.Context, leaveEvent model.LeaveEventRequest) (string, error) {
	// process the data and create new event
	event := database.Event{
		EventID:     leaveEvent.EventID,
		EventTime:   time.Now(),
		EventType:   "RemoveMember",
		Title:       leaveEvent.Title,
		Members:     leaveEvent.NewMembers,
		Description: leaveEvent.Description,
		StartDate:   leaveEvent.StartDate,
		EndDate:     leaveEvent.EndDate,
		Location:    leaveEvent.Location,
	}

	// save event in database
	ok, err := r.repo.RemoveMember(event)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
