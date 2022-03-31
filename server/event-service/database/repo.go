package database

import (
	"context"
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
	InsertEvent(ctx context.Context, event Event) (string, error)
	CreateEvent(ctx context.Context, event Event) (*model.Event, string, error)
	GetEvents(ctx context.Context, username string) ([]*model.Event, []*model.Event, []*model.Event, error)
	RemoveEvent(ctx context.Context, event Event) (string, error)
	EditEvent(ctx context.Context, event Event) (string, string, error)
	AttendEvent(ctx context.Context, eventID, username string, left bool) ([]string, error)
	InsertUserData(ctx context.Context, userData model.UserDataRequest) (*model.UserData, error)
	CheckUserData(ctx context.Context, username string) (*model.UserData, error)
	LogTime(ctx context.Context, eventID, username string, expired bool, leaveTime *time.Time) (string, error)
	CheckIfAttendedOnce(ctx context.Context, username, eventID string) (bool, error)
	CheckIfCurrentlyAttended(ctx context.Context, username, eventID string) (bool, error)
	SetLeaveTimeAfterEventEnded(ctx context.Context, username string, event *model.Event) error
	DetermineContactPersons(ctx context.Context, username, eventID string) ([]string, error)
	GetEvent(ctx context.Context, eventID string) (*model.Event, error)
}

type Repo struct {
	eventCollection    *mongo.Collection
	userDataCollection *mongo.Collection
	timeCollection     *mongo.Collection
}

// NewRepo : create new Repo which includes all necessary database collections
func NewRepo() (Repository, error) {
	client, err := dbConnect()
	if err != nil {
		return nil, err
	}

	db := client.Database("event-service")

	return &Repo{
		eventCollection:    db.Collection("event-events"),
		userDataCollection: db.Collection("user-data"),
		timeCollection:     db.Collection("time"),
	}, nil
}

// InsertEvent for given event
func (repo *Repo) InsertEvent(ctx context.Context, event Event) (string, error) {
	inserted, err := repo.eventCollection.InsertOne(ctx, event)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), err
}

// CreateEvent for event
func (repo *Repo) CreateEvent(ctx context.Context, event Event) (*model.Event, string, error) {
	currentTime := time.Now().UTC()

	// insert event and get insertedID
	insertedID, err := repo.InsertEvent(ctx, event)
	if err != nil {
		return nil, "", err
	}

	// create new event
	eventModel := &model.Event{
		ID:          insertedID,
		Title:       event.Title,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Host:        event.Host,
		Location:    event.Location,
		Description: event.Description,
		Subscribers: event.Subscribers,
		Attendees:   event.Attendees,
	}

	// check which type the created event is: upcoming or ongoing
	startTime, endTime, onlyCheckDate, err := formatDate(eventModel.StartDate, eventModel.EndDate)
	if err != nil {
		return nil, "", err
	}

	// check if event is ongoing
	if !onlyCheckDate && currentTime.Before(endTime) && currentTime.After(startTime) || onlyCheckDate && currentTime.After(startTime) && currentTime.Before(endTime) {
		return eventModel, "ongoing", nil
		// check if event is in the past
	} else if currentTime.After(endTime) {
		return eventModel, "past", nil
		// chcek if event is upcoming
	} else if currentTime.Before(startTime) {
		return eventModel, "upcoming", nil
	}

	return nil, "", errors.New("event is not ongoing, upcoming nor in the past")
}

// GetEvents in three parts: currentEvents, pastEvents and upcomingEvents by username
func (repo *Repo) GetEvents(ctx context.Context, username string) ([]*model.Event, []*model.Event, []*model.Event, error) {
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
			Subscribers: event.Subscribers,
			Host:        event.Host,
			Attendees:   event.Attendees,
		})
	}

	max := int64(1)
	for _, eventModel := range events {
		// Sort event_time and get one element which will be the most recent edited event in relation to liked, unliked and description
		opts.Limit = &max
		cursor, err = repo.eventCollection.Find(ctx, bson.D{
			{"id", eventModel.ID},
			{"event_type", bson.D{
				{"$in", bson.A{"EditEvent", "SubscribeEvent", "AttendEvent"}},
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
			eventModel.Subscribers = event.Subscribers
			eventModel.EndDate = event.EndDate
			eventModel.StartDate = event.StartDate
			eventModel.Location = event.Location
			eventModel.Attendees = event.Attendees
		}

		// Add to correct list in relation to the arrive and endDate of event
		startTime, endTime, onlyCheckDate, err := formatDate(eventModel.StartDate, eventModel.EndDate)
		if err != nil {
			return nil, nil, nil, err
		}

		// check if user is currently attending an ongoing event
		if !onlyCheckDate && currentTime.Before(endTime) && currentTime.After(startTime) || onlyCheckDate && currentTime.After(startTime) && currentTime.Before(endTime) {
			currentlyAttended, err := repo.CheckIfCurrentlyAttended(ctx, username, eventModel.ID)
			if err != nil {
				return nil, nil, nil, err
			}

			eventModel.CurrentlyAttended = &currentlyAttended

			ongoingEvents = append(ongoingEvents, eventModel)
			// check if user is attending an past event and set leave time to endtime of event
		} else if currentTime.After(endTime) {
			attendedOnce, err := repo.CheckIfAttendedOnce(ctx, username, eventModel.ID)
			if err != nil {
				return nil, nil, nil, err
			}

			// only show past events which the user is a host of or attended
			if attendedOnce || eventModel.Host == username {
				// check if user didn't leave the event
				currentlyAttended, err := repo.CheckIfCurrentlyAttended(ctx, username, eventModel.ID)
				if err != nil {
					return nil, nil, nil, err
				}

				// if they're still attending, after event ended set leave time -> endtime of event
				if currentlyAttended {
					err = repo.SetLeaveTimeAfterEventEnded(ctx, username, eventModel)
					if err != nil {
						return nil, nil, nil, err
					}
				}

				pastEvents = append(pastEvents, eventModel)
			}
			// check if event is upcoming
		} else if currentTime.Before(startTime) {
			upcomingEvents = append(upcomingEvents, eventModel)
		}
	}

	return upcomingEvents, ongoingEvents, pastEvents, nil
}

// GetEvent for given event id
func (repo *Repo) GetEvent(ctx context.Context, eventID string) (*model.Event, error) {
	objID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return nil, err
	}

	var event Event
	err = repo.eventCollection.FindOne(ctx, bson.D{{"_id", objID}}).Decode(&event)
	if err != nil {
		return nil, err
	}

	eventModel := &model.Event{
		ID:          event.ID.Hex(),
		Title:       event.Title,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
		Description: event.Description,
		Subscribers: event.Subscribers,
		Host:        event.Host,
		Attendees:   event.Attendees,
	}

	max := int64(1)
	// Sort event_time and get one element which will be the most recent edited event in relation to liked, unliked and description
	opts := options.Find()
	opts.SetSort(bson.D{{"event_time", -1}})
	opts.Limit = &max
	cursor, err := repo.eventCollection.Find(ctx, bson.D{
		{"id", eventID},
		{"event_type", bson.D{
			{"$in", bson.A{"EditEvent", "SubscribeEvent", "AttendEvent"}},
		}},
	}, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var event Event
		if err = cursor.Decode(&event); err != nil {
			return nil, err
		}

		// Add editable data
		eventModel.Title = event.Title
		eventModel.Description = event.Description
		eventModel.Subscribers = event.Subscribers
		eventModel.EndDate = event.EndDate
		eventModel.StartDate = event.StartDate
		eventModel.Location = event.Location
		eventModel.Attendees = event.Attendees
	}

	return eventModel, nil
}

// formatDate for given input
func formatDate(startDate string, endDate string) (time.Time, time.Time, bool, error) {
	// event has a timestamp
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
		// event has no timestamp
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

// RemoveEvent and all its created events
func (repo *Repo) RemoveEvent(ctx context.Context, event Event) (string, error) {
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
			{"$in", bson.A{"EditEvent", "SubscribeEvent", "AttendEvent"}},
		}},
	})
	if err != nil {
		return "failed", err
	}

	// delete all logTimes
	_, err = repo.timeCollection.DeleteMany(ctx, bson.D{
		{"id", event.EventID},
	})
	if err != nil {
		return "failed", err
	}

	// new current event events
	_, err = repo.InsertEvent(ctx, event)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

// EditEvent with given event
func (repo *Repo) EditEvent(ctx context.Context, event Event) (string, string, error) {
	_, err := repo.InsertEvent(ctx, event)
	if err != nil {
		return "failed", "", err
	}

	currentTime := time.Now().UTC()

	// format date to go time lib
	startTime, endTime, onlyCheckDate, err := formatDate(event.StartDate, event.EndDate)
	if err != nil {
		return "failure", "", err
	}

	// check edited events time (when edited, it may change the category)
	if !onlyCheckDate && currentTime.Before(endTime) && currentTime.After(startTime) || onlyCheckDate && currentTime.After(startTime) && currentTime.Before(endTime) {
		return "success", "ongoing", nil
	} else if currentTime.After(endTime) {
		return "success", "past", nil
	} else if currentTime.Before(startTime) {
		return "success", "upcoming", nil
	}

	return "failure", "", errors.New("event is not ongoing, upcoming nor in the past")
}

// AttendEvent with given eventId, username and if attending or leaving condition
func (repo *Repo) AttendEvent(ctx context.Context, eventID, username string, left bool) ([]string, error) {
	// get most recent edited event
	updatedEvent, err := repo.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// update attendees list
	if !left {
		updatedEvent.Attendees = append(updatedEvent.Attendees, username)
	} else {
		index := -1
		for i, attendee := range updatedEvent.Attendees {
			if attendee == username {
				index = i
				break
			}
		}

		if index == -1 {
			return nil, errors.New("username not found in attendees list")
		}

		updatedEvent.Attendees = append(updatedEvent.Attendees[:index], updatedEvent.Attendees[index+1:]...)
	}

	// process the data and create new event
	event := Event{
		EventID:     updatedEvent.ID,
		EventTime:   time.Now(),
		EventType:   "AttendEvent",
		Title:       updatedEvent.Title,
		Subscribers: updatedEvent.Subscribers,
		Description: updatedEvent.Description,
		StartDate:   updatedEvent.StartDate,
		EndDate:     updatedEvent.EndDate,
		Location:    updatedEvent.Location,
		Attendees:   updatedEvent.Attendees,
	}

	// get utc time now because all events are timezone utc
	currentTime := time.Now().UTC()

	// format endTime because it is needed for error handling
	_, endTime, _, _ := formatDate(event.StartDate, event.EndDate)

	// check if current time is before endtime of event
	if currentTime.Before(endTime) {
		// if user is attending event, a new event will be inserted to db
		if !left {
			_, err := repo.InsertEvent(ctx, event)
			if err != nil {
				return nil, err
			}
		}

		// log the time of the users attendance
		_, err := repo.LogTime(ctx, event.EventID, username, false, nil)
		if err != nil {
			return nil, err
		}

		return updatedEvent.Attendees, nil
	}

	return nil, errors.New("event is expired")
}

// InsertUserData with given userdata
func (repo *Repo) InsertUserData(ctx context.Context, userData model.UserDataRequest) (*model.UserData, error) {
	_, err := repo.userDataCollection.InsertOne(ctx, userData)
	if err != nil {
		return nil, err
	}

	return &model.UserData{
		Username:    userData.Username,
		Firstname:   userData.Firstname,
		Lastname:    userData.Lastname,
		Street:      userData.Street,
		Housenumber: userData.Housenumber,
		Postalcode:  userData.Postalcode,
		City:        userData.City,
		Email:       userData.Email,
	}, nil
}

// CheckUserData if userdata exists in db
func (repo *Repo) CheckUserData(ctx context.Context, username string) (*model.UserData, error) {
	var result model.UserData
	err := repo.userDataCollection.FindOne(ctx, bson.D{{"username", username}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

// LogTime of user attendance
func (repo *Repo) LogTime(ctx context.Context, eventID, username string, expired bool, leaveTime *time.Time) (string, error) {
	// check if user attended
	var log LogTime
	err := repo.timeCollection.FindOne(ctx, bson.D{
		{"username", username},
		{"id", eventID},
		{"leave", bson.D{
			{"$exists", false},
		}},
	}).Decode(&log)
	if err != nil {
		// if user didn't attend, create a new document for attendance
		if err == mongo.ErrNoDocuments {
			logTime := LogTime{
				EventID:  eventID,
				Username: username,
				Arrive:   time.Now(),
			}

			_, err = repo.timeCollection.InsertOne(ctx, logTime)
			if err != nil {
				return "", err
			}

			return "success", nil
		}
		return "", err
	}

	// update an attendance document with currentTime or leaveTime (endDate of event)
	filter := bson.D{{"_id", log.ID}}
	var update bson.D
	if expired {
		update = bson.D{{"$set", bson.D{{"leave", *leaveTime}}}}
	} else {
		update = bson.D{{"$set", bson.D{{"leave", time.Now()}}}}
	}

	result, err := repo.timeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	if result.MatchedCount == 0 {
		text := "no event with id " + log.EventID + " found"
		return "", errors.New(text)
	}

	return "success", nil
}

// CheckIfAttendedOnce : check if a user attended the event with given eventID
func (repo *Repo) CheckIfAttendedOnce(ctx context.Context, username, eventID string) (bool, error) {
	var log LogTime
	err := repo.timeCollection.FindOne(ctx, bson.D{
		{"username", username},
		{"id", eventID}}).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// CheckIfCurrentlyAttended : check if user is currently attending an event
func (repo *Repo) CheckIfCurrentlyAttended(ctx context.Context, username, eventID string) (bool, error) {
	var log LogTime
	err := repo.timeCollection.FindOne(ctx, bson.D{
		{"username", username},
		{"id", eventID},
		{"leave", bson.D{
			{"$exists", false},
		}},
	}).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// SetLeaveTimeAfterEventEnded : if an event was expired the leave time of an open attendance will be set to the end time of event
func (repo *Repo) SetLeaveTimeAfterEventEnded(ctx context.Context, username string, event *model.Event) error {
	_, leaveTime, _, err := formatDate(event.StartDate, event.EndDate)
	if err != nil {
		return err
	}

	_, err = repo.LogTime(ctx, event.ID, username, true, &leaveTime)

	if err != nil {
		return err
	}

	return nil
}

// DetermineContactPersons : get contact persons of an event for given username and eventID
func (repo *Repo) DetermineContactPersons(ctx context.Context, username, eventID string) ([]string, error) {
	contactPersons := make([]string, 0)

	type interval struct {
		arrive time.Time
		leave  time.Time
	}

	logTimes := make([]interval, 0)

	cursor, err := repo.timeCollection.Find(ctx, bson.D{
		{"id", eventID},
		{"username", username},
	})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var logTime LogTime
		if err = cursor.Decode(&logTime); err != nil {
			return nil, err
		}

		logTimes = append(logTimes, interval{
			arrive: logTime.Arrive,
			leave:  logTime.Leave,
		})
	}

	for _, logTime := range logTimes {
		cursor, err := repo.timeCollection.Find(ctx, bson.D{
			{"id", eventID},
			{"username", bson.D{
				{"$ne", username},
			}},
			{"$or", []interface{}{
				bson.D{
					{"$and", []interface{}{
						bson.D{{"arrive", bson.D{
							{"$lte", logTime.arrive},
						}}},
						bson.D{{"leave", bson.D{
							{"$gte", logTime.arrive},
						}}},
					}},
				},
				bson.D{
					{"$and", []interface{}{
						bson.D{{"arrive", bson.D{
							{"$gte", logTime.arrive},
							{"$lte", logTime.leave},
						}}},
					}},
				},
			}},
		})
		if err != nil {
			return nil, err
		}

		for cursor.Next(ctx) {
			var logTime LogTime
			if err = cursor.Decode(&logTime); err != nil {
				return nil, err
			}

			contactPersons = append(contactPersons, logTime.Username)
		}
	}

	removeDuplicatesStr := func(strSlice []string) []string {
		allKeys := make(map[string]bool)
		list := []string{}
		for _, item := range strSlice {
			if _, value := allKeys[item]; !value {
				allKeys[item] = true
				list = append(list, item)
			}
		}
		return list
	}

	return removeDuplicatesStr(contactPersons), nil
}
