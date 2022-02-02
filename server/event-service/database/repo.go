package database

import (
	"errors"
	"event-service/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type Repository interface {
	InsertEvent(event Event) (string, error)
	CreateEvent(event Event) (*model.Event, string, error)
	GetEvents() ([]*model.Event, []*model.Event, []*model.Event, error)
	RemoveEvent(event Event) (string, error)
	EditEvent(event Event) (string, error)
	AddedMember(event Event) (string, error)
	RemoveMember(event Event) (string, error)
}

type Repo struct {
	eventCollection    *mongo.Collection
	userDataCollection *mongo.Collection
}

func NewRepo() (Repository, error) {
	client, err := dbConnect()
	if err != nil {
		return nil, err
	}

	db := client.Database("event-service")

	return &Repo{
		eventCollection:    db.Collection("event-events"),
		userDataCollection: db.Collection("user-data"),
	}, nil
}

func (repo *Repo) InsertEvent(event Event) (string, error) {
	inserted, err := repo.eventCollection.InsertOne(ctx, event)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), err
}

func (repo *Repo) CreateEvent(event Event) (*model.Event, string, error) {
	currentTime := time.Now().UTC()

	insertedID, err := repo.InsertEvent(event)
	if err != nil {
		return nil, "", err
	}

	eventModel := &model.Event{
		ID:          insertedID,
		Title:       event.Title,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Host:        event.Host,
		Location:    event.Location,
		Description: event.Description,
		Members:     event.Members,
	}

	startTime, endTime, onlyCheckDate, err := formatDate(eventModel.StartDate, eventModel.EndDate)
	if err != nil {
		return nil, "", err
	}

	if !onlyCheckDate && currentTime.Before(endTime) && currentTime.After(startTime) || onlyCheckDate && currentTime.After(startTime) && currentTime.Before(endTime) {
		return eventModel, "ongoing", nil
	} else if currentTime.After(endTime) {
		return eventModel, "past", nil
	} else if currentTime.Before(startTime) {
		return eventModel, "upcoming", nil
	}

	return nil, "", errors.New("event is whether ongoing, upcoming nor in the past")
}

func (repo *Repo) GetEvents() ([]*model.Event, []*model.Event, []*model.Event, error) {
	events, upcomingEvents, ongoingEvents, pastEvents := make([]*model.Event, 0), make([]*model.Event, 0), make([]*model.Event, 0), make([]*model.Event, 0)
	currentTime := time.Now().UTC()

	// sort eventModel-events by descending eventModel-time (the newest first) and set fetch limit
	opts := options.Find()
	opts.SetSort(bson.D{{"event_time", -1}})

	// get all eventModel events with event_type = "CreateEvent" sorted by event_time
	cursor, err := repo.eventCollection.Find(ctx, bson.D{
		{"event_type", "CreateEvent"},
	}, opts)
	if err != nil {
		return nil, nil, nil, err
	}

	for cursor.Next(ctx) {
		var event Event
		if err = cursor.Decode(&event); err != nil {
			return nil, nil, nil, err
		}

		// add new eventModel to output for getEvents
		events = append(events, &model.Event{
			ID:          event.ID.Hex(),
			Title:       event.Title,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			Location:    event.Location,
			Description: event.Description,
			Members:     event.Members,
			Host:        event.Host,
		})
	}

	max := int64(1)
	for _, eventModel := range events {
		// Sort event_time and get one element which will be the most recent edited event in relation to liked, unliked and description
		opts.Limit = &max
		cursor, err = repo.eventCollection.Find(ctx, bson.D{
			{"id", eventModel.ID},
			{"event_type", bson.D{
				{"$in", bson.A{"EditEvent", "AddedMember", "RemoveMember"}},
			}},
		}, opts)
		if err != nil {
			return nil, nil, nil, err
		}

		for cursor.Next(ctx) {
			var event Event
			if err = cursor.Decode(&event); err != nil {
				return nil, nil, nil, err
			}

			// Add editable data
			eventModel.Title = event.Title
			eventModel.Description = event.Description
			eventModel.Members = event.Members
			eventModel.EndDate = event.EndDate
			eventModel.StartDate = event.StartDate
			eventModel.Location = event.Location
		}

		// Add to correct list in relation to the start and endDate of event
		startTime, endTime, onlyCheckDate, err := formatDate(eventModel.StartDate, eventModel.EndDate)
		if err != nil {
			return nil, nil, nil, err
		}

		if !onlyCheckDate && currentTime.Before(endTime) && currentTime.After(startTime) || onlyCheckDate && currentTime.After(startTime) && currentTime.Before(endTime) {
			ongoingEvents = append(ongoingEvents, eventModel)
		} else if currentTime.After(endTime) {
			pastEvents = append(pastEvents, eventModel)
		} else if currentTime.Before(startTime) {
			upcomingEvents = append(upcomingEvents, eventModel)
		}
	}

	return upcomingEvents, ongoingEvents, pastEvents, nil
}

func formatDate(startDate string, endDate string) (time.Time, time.Time, bool, error) {
	if strings.HasSuffix(startDate, "M") && strings.HasSuffix(endDate, "M") {
		startTime, err := time.Parse("1/2/06, 3:04 PM", startDate)
		if err != nil {
			return time.Time{}, time.Time{}, false, err
		}

		endTime, err := time.Parse("1/2/06, 3:04 PM", endDate)
		if err != nil {
			return time.Time{}, time.Time{}, false, err
		}

		return startTime, endTime, false, nil
	} else {
		startTime, err := time.Parse("Monday, January 2, 2006", startDate)
		if err != nil {
			return time.Time{}, time.Time{}, false, err
		}

		endTime, err := time.Parse("Monday, January 2, 2006", endDate)
		if err != nil {
			return time.Time{}, time.Time{}, false, err
		}

		endTime = endTime.Add(time.Hour * 24)

		return startTime, endTime, true, nil
	}
}

func (repo *Repo) RemoveEvent(event Event) (string, error) {
	// convert hex-string into primitive.objectID
	objID, err := primitive.ObjectIDFromHex(event.EventID)
	if err != nil {
		return "failed", err
	}

	// delete that one CreateEvent-Event
	_, err = repo.eventCollection.DeleteOne(ctx, bson.D{
		{"_id", objID},
		{"event_type", "CreateEvent"},
	})
	if err != nil {
		return "failed", err
	}

	// delete all other events
	_, err = repo.eventCollection.DeleteMany(ctx, bson.D{
		{"id", event.EventID},
		{"event_type", bson.D{
			{"$in", bson.A{"EditEvent", "AddedMember", "RemoveMember"}},
		}},
	})
	if err != nil {
		return "failed", err
	}

	// new current event events
	_, err = repo.InsertEvent(event)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *Repo) EditEvent(event Event) (string, error) {
	_, err := repo.InsertEvent(event)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *Repo) AddedMember(event Event) (string, error) {
	_, err := repo.InsertEvent(event)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *Repo) RemoveMember(event Event) (string, error) {
	_, err := repo.InsertEvent(event)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}
