package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-service/graph/model"
)

type PostEvent struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	EventTime   string             `bson:"event_time,omitempty"`
	EventType   string             `bson:"event_type,omitempty"`
	PostID      string             `bson:"id,omitempty"`
	Username    string             `bson:"username,omitempty"`
	Description string             `bson:"description,omitempty"`
	FileID      string             `bson:"fileID,omitempty"`
	LikedBy     []string           `bson:"liked_by,omitempty"`
	Comments    []string           `bson:"comments,omitempty"`
}

type CommentEvent struct {
	Comment *model.Comment `json:"comment"`
	Post    *model.Post    `json:"post"`
}