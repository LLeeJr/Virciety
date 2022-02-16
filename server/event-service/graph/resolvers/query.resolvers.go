package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"event-service/graph/generated"
	"event-service/graph/model"
)

func (r *queryResolver) GetEvents(ctx context.Context) (*model.GetEventsResponse, error) {
	upcomingEvents, ongoingEvents, pastEvents, err := r.repo.GetEvents()
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

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
