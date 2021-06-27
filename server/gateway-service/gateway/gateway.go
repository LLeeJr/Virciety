package gateway

import "context"

type DMEvent struct {
	EventTime string `json:"eventTime"`
	From string `json:"from"`
	To string `json:"to"`
	Time string `json:"time"`
	Payload DM `json:"payload"`
}

type DM struct {
	ID string `json:"id"`
	Msg string `json:"msg"`
}

type Repository interface {
	CreateDMEvent(ctx context.Context, event DMEvent) error
	GetDM(ctx context.Context, id string) (string, error)
}
