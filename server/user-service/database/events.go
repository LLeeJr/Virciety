package database

import "time"

type UserEvent struct {
	EventType string    `json:"eventType"`
	EventTime time.Time `json:"eventTime"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Follows   []string  `json:"follows"`
	Followers []string  `json:"followers"`
}

type FollowEvent struct {
	EventType string    `json:"eventType"`
	EventTime time.Time `json:"eventTime"`
	Follows   string    `json:"follows"`
	ID        string    `json:"id"`
}