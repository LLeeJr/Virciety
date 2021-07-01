package database

import (
	"context"
	"database/sql"
	"errors"
	"notifs-service/graph/model"
)

type NotifEvent struct {
	EventTime string `json:"eventTime"`
	EventType string `json:"eventType"`
	NotifId string `json:"id"`
	Receiver string `json:"receiver"`
	Text string `json:"text"`
}

type Repository interface {
	CreateNotif(ctx context.Context, notifEvent NotifEvent) (*model.Notif, error)
	GetNotifsByReceiver(ctx context.Context, receiver string) ([]*model.Notif, error)
}

type repo struct {
	DB *sql.DB
	NotifEvents []*NotifEvent
}

func (r repo) CreateNotif(ctx context.Context, notifEvent NotifEvent) (*model.Notif, error) {
	r.NotifEvents = append(r.NotifEvents, &notifEvent)

	query := `INSERT INTO db.public.notifs ("eventTime", "eventType", "NotifId", "Receiver", "Text")
              VALUES ($1, $2, $3, $4, $5)`

	_, err := r.DB.ExecContext(ctx, query, notifEvent.EventTime, notifEvent.EventType,
		notifEvent.NotifId, notifEvent.Receiver, notifEvent.Text)
	if err != nil {
		return nil, err
	}

	return &model.Notif{
		ID:       notifEvent.NotifId,
		Receiver: notifEvent.Receiver,
		Text:     notifEvent.Text,
		Event:    notifEvent.EventType,
	}, nil

}

func (r repo) GetNotifsByReceiver(ctx context.Context, receiver string) ([]*model.Notif, error) {
	notifs := make([]*model.Notif, 0)

	query := `SELECT * FROM db.public.notifs WHERE "Receiver" = $1`

	rows, _ := r.DB.QueryContext(ctx, query, receiver)
	for rows.Next() {
		var notifEvent NotifEvent
		err := rows.Scan(&notifEvent.EventTime, &notifEvent.EventType,
			&notifEvent.NotifId, &notifEvent.Receiver, &notifEvent.Text)
		if err != nil {
			continue
		}
		notif := &model.Notif{
			ID:       notifEvent.NotifId,
			Receiver: notifEvent.Receiver,
			Text:     notifEvent.Text,
			Event:    notifEvent.EventType,
		}
		notifs = append(notifs, notif)
	}

	if len(notifs) == 0 {
		return nil, errors.New("no notifs available")
	}

	return notifs, nil
}

func NewRepo(db *sql.DB) (Repository, error) {
	return &repo{
		DB:       db,
		NotifEvents: make([]*NotifEvent, 0),
	}, nil
}