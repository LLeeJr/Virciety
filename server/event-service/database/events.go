package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	EventID     string             `bson:"id,omitempty"`
	EventTime   time.Time          `bson:"event_time,omitempty"`
	EventType   string             `bson:"event_type,omitempty"`
	Members     []string           `bson:"members,omitempty"`
	Host        string             `bson:"host,omitempty"`
	Description string             `bson:"description,omitempty"`
	StartDate   string             `bson:"startDate,omitempty"`
	EndDate     string             `bson:"endDate,omitempty"`
	Location    string             `bson:"location,omitempty"`
}
