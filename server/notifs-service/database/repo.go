package database

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"notifs-service/graph/model"
	"time"
)

type Message struct {
	Observers map[string]struct{
		Message chan *model.Notif
	}
}


type NotifEvent struct {
	EventTime time.Time `json:"eventtime"`
	EventType string    `json:"eventtype"`
	Receiver  []string  `json:"receiver"`
	Text      string    `json:"text"`
}

type ChatEvent struct {
	EventTime time.Time `json:"eventTime"`
	From      string    `json:"from"`
	Msg       string    `json:"msg"`
	RoomID    string    `json:"roomId"`
	RoomName  string    `json:"roomName"`
	Receivers []string  `json:"receivers"`
}

type Repository interface {
	CreateNotif(ctx context.Context, notifEvent NotifEvent) (*model.Notif, error)
	GetNotifsByReceiver(ctx context.Context, receiver string) ([]*model.Notif, error)
	CreateDmNotifFromConsumer(body []byte) error
	GetSubscriptions() map[string]*Message
	AddSubscription(name string, subscription *Message)
}

type repo struct {
	notifCollection *mongo.Collection
	NotifEvents []*NotifEvent
	Subscriptions map[string]*Message
}

func (r repo) AddSubscription(name string, subscription *Message) {
	r.Subscriptions[name] = subscription
}

func (r repo) GetSubscriptions() map[string]*Message {
	return r.Subscriptions
}

func (r repo) CreateDmNotifFromConsumer(data []byte) error {
	var s *ChatEvent
	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println("err", err)
		return err
	}

	notifText := fmt.Sprintf("New message from %s in room %s", s.From, s.RoomName)
	notifEvent := NotifEvent{
		EventTime: s.EventTime,
		EventType: "New DM",
		Receiver:  s.Receivers,
		Text:      notifText,
	}

	insertedId, err := r.InsertNotifEvent(context.Background(), notifEvent)
	if err != nil {
		return err
	}
	log.Println("new notif: ", insertedId)
	r.NotifEvents = append(r.NotifEvents, &notifEvent)

	for _, receiver := range s.Receivers {
		subscription := r.Subscriptions[receiver]
		if subscription != nil {
			notif := &model.Notif{
				ID:       insertedId,
				Receiver: notifEvent.Receiver,
				Text:     notifEvent.Text,
				Event:    notifEvent.EventType,
			}
			for _, observer := range subscription.Observers {
				observer.Message <- notif
			}
		}
	}

	return nil
}

func (r *repo) InsertNotifEvent(ctx context.Context, notifEvent NotifEvent) (string, error) {
	inserted, err := r.notifCollection.InsertOne(ctx, notifEvent)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r repo) CreateNotif(ctx context.Context, notifEvent NotifEvent) (*model.Notif, error) {
	r.NotifEvents = append(r.NotifEvents, &notifEvent)

	insertedId, err := r.InsertNotifEvent(ctx, notifEvent)
	if err != nil {
		return nil, err
	}

	return &model.Notif{
		ID:       insertedId,
		Receiver: notifEvent.Receiver,
		Text:     notifEvent.Text,
		Event:    notifEvent.EventType,
	}, nil

}

func (r repo) GetNotifsByReceiver(ctx context.Context, receiver string) ([]*model.Notif, error) {

	type Notif struct {
		ID        primitive.ObjectID `bson:"_id"`
		EventTime time.Time          `bson:"eventtime"`
		EventType string             `bson:"eventtype"`
		Receiver  []string           `bson:"receiver"`
		Text      string             `bson:"text"`
	}

	var result []*Notif

	opts := options.Find().SetSort(bson.D{{"eventtime", -1}}).SetLimit(10)
	cursor, err := r.notifCollection.Find(
		ctx,
		bson.D{
			{"receiver", receiver},
		}, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	notifs := make([]*model.Notif, 0)
	for _, notif := range result {
		notifs = append(notifs, &model.Notif{
			ID:       notif.ID.Hex(),
			Receiver: notif.Receiver,
			Text:     notif.Text,
			Event:    notif.EventType,
		})
	}

	return notifs, nil
}

func NewRepo() (Repository, error) {
	client, err := Connect()
	if err != nil {
		return nil, err
	}

	db := client.Database("notif-service")

	return &repo{
		notifCollection: db.Collection("notif-events"),
		NotifEvents: make([]*NotifEvent, 0),
		Subscriptions: map[string]*Message{},
	}, nil
}