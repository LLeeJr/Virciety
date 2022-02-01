package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"event-service/graph/generated"
	"event-service/graph/model"
	"strings"
	"time"
)

func (r *queryResolver) GetEvents(ctx context.Context) (*model.GetEventsResponse, error) {
	upcomingEvents, ongoingEvents, pastEvents := make([]*model.Event, 0), make([]*model.Event, 0), make([]*model.Event, 0)

	currentTime := time.Now().UTC()
	allEvents, _ := r.repo.GetEvents()

	for _, event := range allEvents {
		var startTime, endTime time.Time
		var err error
		var onlyCheckDate bool

		if strings.HasSuffix(event.StartDate, "M") && strings.HasSuffix(event.EndDate, "M") {
			onlyCheckDate = false
			startTime, err = time.Parse("1/2/06, 3:04 PM", event.StartDate)
			if err != nil {
				return nil, err
			}

			endTime, err = time.Parse("1/2/06, 3:04 PM", event.EndDate)
			if err != nil {
				return nil, err
			}
		} else {
			onlyCheckDate = true
			startTime, err = time.Parse("Monday, January 2, 2006", event.StartDate)
			if err != nil {
				return nil, err
			}

			endTime, err = time.Parse("Monday, January 2, 2006", event.EndDate)
			if err != nil {
				return nil, err
			}

			endTime = endTime.Add(time.Hour * 24)
		}

		if !onlyCheckDate && currentTime.Before(endTime) && currentTime.After(startTime) || onlyCheckDate && currentTime.After(startTime) && currentTime.Before(endTime) {
			ongoingEvents = append(ongoingEvents, event)
		} else if currentTime.After(endTime) {
			pastEvents = append(pastEvents, event)
		} else if currentTime.Before(startTime) {
			upcomingEvents = append(upcomingEvents, event)
		}
	}

	return &model.GetEventsResponse{
		UpcomingEvents: upcomingEvents,
		OngoingEvents:  ongoingEvents,
		PastEvents:     pastEvents,
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
