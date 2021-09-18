package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostEvent struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	EventTime   string             `bson:"event_time,omitempty"`
	EventType   string             `bson:"event_type,omitempty"`
	PostID      string             `bson:"id,omitempty"` // replace with id from mongodb
	Username    string             `bson:"username,omitempty"`
	Description string             `bson:"description,omitempty"`
	FileID      string             `bson:"fileID,omitempty"`
	LikedBy     []string           `bson:"liked_by,omitempty"`
	Comments    []string           `bson:"comments,omitempty"`
}

type CommentEvent struct {
	EventTime   string   `json:"event_time"`
	EventType   string   `json:"event_type"`
	CommentID   string   `json:"id"`
	Username    string   `json:"username"`
	Description string   `json:"description"`
	LikedBy     []string `json:"liked_by"`
	PostID      string   `json:"post_id"`
}
