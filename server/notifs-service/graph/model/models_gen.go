// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Map struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Notif struct {
	ID        string    `json:"id"`
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Read      bool      `json:"read"`
	Receiver  string    `json:"receiver"`
	Text      string    `json:"text"`
	Params    []*Map    `json:"params"`
	Route     string    `json:"route"`
}
