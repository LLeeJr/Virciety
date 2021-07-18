package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"notifs-service/graph/model"
	"strings"
	"time"
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
	CreateDmNotifFromConsumer(body []byte)
}

type repo struct {
	DB *sql.DB
	NotifEvents []*NotifEvent
}

func (r repo) CreateDmNotifFromConsumer(data []byte) {
	var s map[string]string
	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println("err", err)
		return
	}

	id := s["id"]
	slices := strings.Split(id, "__")
	notifId := fmt.Sprintf("%s__%s__%s", s["eventType"], time.Now().Format("2006-01-02 15:04:05"), slices[2])
	notifText := fmt.Sprintf("Got a new DM from %s", slices[0])
	notifEvent := NotifEvent{
		EventTime: s["eventTime"],
		EventType: s["eventType"],
		NotifId:   notifId,
		Receiver:  slices[2],
		Text:      notifText,
	}

	r.NotifEvents = append(r.NotifEvents, &notifEvent)
	r.InsertNotifEvent(notifEvent)
}

func (r *repo) InsertNotifEvent(notifEvent NotifEvent) (err error) {
	query := `INSERT INTO db.public.notifs ("eventTime", "eventType", "NotifId", "Receiver", "Text")
              VALUES ($1, $2, $3, $4, $5)`

	_, err = r.DB.Exec(query, notifEvent.EventTime, notifEvent.EventType,
		notifEvent.NotifId, notifEvent.Receiver, notifEvent.Text)
	if err != nil {
		return err
	}

	return
}

func (r repo) CreateNotif(ctx context.Context, notifEvent NotifEvent) (*model.Notif, error) {
	r.NotifEvents = append(r.NotifEvents, &notifEvent)

	r.InsertNotifEvent(notifEvent)

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