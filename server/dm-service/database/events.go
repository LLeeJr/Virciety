package database

import (
	"time"
)

type DmEvent struct {
	ChatroomId string `json:"id"`
	CreatedAt  time.Time `json:"eventTime"`
	CreatedBy  string `json:"to"`
	EventType  string    `json:"eventType"`
	Msg        string `json:"msg"`
}

type ChatroomEvent struct {
	EventType  string   `json:"eventType"`
	Member     []string `json:"member"`
	Name       string   `json:"name"`
	Owner      string   `json:"owner"`
	MemberSize int      `json:"membersize"`
	IsDirect   *bool    `json:"isdirect"`
}