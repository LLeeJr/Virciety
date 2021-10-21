package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommentEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	EventTime string             `bson:"event_time,omitempty"`
	EventType string             `bson:"event_type,omitempty"`
	CommentID string             `bson:"id,omitempty"`
	PostID    string             `bson:"post_id,omitempty"`
	Comment   string             `bson:"comment,omitempty"`
	CreatedBy string             `bson:"created_by,omitempty"`
}
