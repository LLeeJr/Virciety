package database

import (
	"time"
)

type UserEvent struct {
	EventType string    `json:"eventType"`
	EventTime time.Time `json:"eventTime"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Follows   []string  `json:"follows"`
	Followers []string  `json:"followers"`
	FileId    string    `json:"fileId"`
}

type ProfilePictureEvent struct {
	EventType string    `json:"eventType"`
	EventTime time.Time `json:"eventTime"`
	FileId    string    `json:"fileID"`
	Username  string    `json:"username"`
}

type FollowEvent struct {
	EventType   string    `json:"event_type"`
	EventTime   time.Time `json:"event_time"`
	Username    string    `json:"username"`
	NewFollower string    `json:"new_follower"`
}