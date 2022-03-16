package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"event-service/graph/generated"
	"event-service/graph/model"
	message_queue "event-service/message-queue"
)

func (r *queryResolver) GetEvents(ctx context.Context, username *string) (*model.GetEventsResponse, error) {
	upcomingEvents, ongoingEvents, pastEvents, err := r.repo.GetEvents(ctx, *username)
	if err != nil {
		return nil, err
	}

	return &model.GetEventsResponse{
		UpcomingEvents: upcomingEvents,
		OngoingEvents:  ongoingEvents,
		PastEvents:     pastEvents,
	}, nil
}

func (r *queryResolver) UserDataExists(ctx context.Context, username *string) (*model.UserData, error) {
	userData, err := r.repo.CheckUserData(ctx, *username)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (r *queryResolver) NotifyHostOfEvent(ctx context.Context, username *string, eventID *string, reportedBy *string) (*bool, error) {
	eventNotification := message_queue.EventNotification{
		EventId:    *eventID,
		EditFlag:   false,
		Message:    "A covid case was reported",
		ReportedBy: *reportedBy,
		Username:   *username,
	}

	r.producerQueue.AddMessageToEvent(eventNotification)

	notified := true

	return &notified, nil
}

func (r *queryResolver) NotifyContactPersons(ctx context.Context, username *string, eventID *string) (*bool, error) {
	notified := true

	contactPersons, err := r.repo.DetermineContactPersons(ctx, *username, *eventID)
	if err != nil {
		return nil, err
	}

	eventNotification := message_queue.EventNotification{
		EventId:  *eventID,
		EditFlag: false,
		Message:  "You were in contact with a person who was tested positive for covid. Please, get yourself tested.",
	}
	for _, person := range contactPersons {
		eventNotification.Username = person
		r.producerQueue.AddMessageToEvent(eventNotification)
	}

	return &notified, nil
}

func (r *queryResolver) GetEvent(ctx context.Context, eventID *string) (*model.Event, error) {
	return r.repo.GetEvent(ctx, *eventID);
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
