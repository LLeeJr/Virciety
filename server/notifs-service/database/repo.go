package database

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"notifs-service/graph/model"
	"time"
)

type Message struct {
	Observers map[string]struct{
		Message chan *model.Notif
	}
}


type NotifEvent struct {
	EventTime time.Time   `json:"eventtime"`
	EventType string      `json:"eventtype"`
	Receiver  string      `json:"receiver"`
	Text      string      `json:"text"`
	Read      bool        `json:"read"`
	Route     string      `json:"route"`
	Params    []*model.Map `json:"params"`
}

type ChatEvent struct {
	EventTime time.Time `json:"eventTime"`
	From      string    `json:"from"`
	Msg       string    `json:"msg"`
	RoomID    string    `json:"roomId"`
	RoomName  string    `json:"roomName"`
	Receivers []string  `json:"receivers"`
}

type Comment struct {
	ID        string `json:"id"`
	PostID    string `json:"postID"`
	Comment   string `json:"comment"`
	CreatedBy string `json:"createdBy"`
	Event     string `json:"event"`
}

type Post struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Username    string   `json:"username"`
}

type CommentEvent struct {
	Comment Comment `json:"comment"`
	Post    Post    `json:"post"`
}

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

type Repository interface {
	GetNotifsByReceiver(ctx context.Context, receiver string) ([]*model.Notif, error)
	CreateDmNotifFromConsumer(body []byte) error
	GetSubscriptions() map[string]*Message
	AddSubscription(name string, subscription *Message)
	GetNotification(ctx context.Context, id string) (*model.Notif, error)
	UpdateNotification(ctx context.Context, id string, status bool) (string, error)
	CreateCommentNotifFromConsumer(body []byte) error
	CreateLikeNotifFromConsumer(body []byte) error
}

type repo struct {
	notifCollection *mongo.Collection
	Subscriptions map[string]*Message
}

type Notif struct {
	ID        primitive.ObjectID `bson:"_id"`
	EventTime time.Time          `bson:"eventtime"`
	EventType string             `bson:"eventtype"`
	Params    []*model.Map       `bson:"params"`
	Receiver  string             `bson:"receiver"`
	Read      bool               `bson:"read"`
	Route     string             `bson:"route"`
	Text      string             `bson:"text"`
}

func (r repo) UpdateNotification(ctx context.Context, id string, status bool) (string, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	query := bson.M{
		"_id": objID,
	}

	update := bson.M{
		"$set": bson.D{{"read", status}},
	}

	_, err = r.notifCollection.UpdateOne(ctx, query, update)
	if err != nil {
		return "", err
	}

	return "success", nil
}

func (r repo) GetNotification(ctx context.Context, id string) (*model.Notif, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result *Notif
	if err := r.notifCollection.FindOne(ctx, bson.D{
		{"_id", objID},
	}).Decode(&result); err != nil {
		return nil, err
	}

	notif := &model.Notif{
		ID:       id,
		Event:    result.EventType,
		Read:     result.Read,
		Receiver: result.Receiver,
		Text:     result.Text,
		Timestamp: result.EventTime,
		Params:   result.Params,
		Route:    result.Route,
	}

	return notif, nil
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
		return err
	}

	m := []*model.Map{
		{
			Key:   "from",
			Value: s.From,
		},
		{
			Key:   "roomName",
			Value: s.RoomName,
		},
		{
			Key:   "roomID",
			Value: s.RoomID,
		},
	}

	notifText := fmt.Sprintf("New message from %s in room %s", s.From, s.RoomName)
	notifEvent := NotifEvent{
		EventTime: s.EventTime,
		EventType: "New DM",
		Params:    m,
		Read:      false,
		Route:     "/chat",
		Text:      notifText,
	}

	for _, receiver := range s.Receivers {

		notifEvent.Receiver = receiver
		insertedId, err := r.InsertNotifEvent(context.Background(), notifEvent)
		if err != nil {
			return err
		}

		subscription := r.Subscriptions[receiver]
		if subscription != nil {
			notif := &model.Notif{
				ID:       insertedId,
				Params:   notifEvent.Params,
				Receiver: notifEvent.Receiver,
				Read:     notifEvent.Read,
				Route:    notifEvent.Route,
				Text:     notifEvent.Text,
				Timestamp: notifEvent.EventTime,
				Event:    notifEvent.EventType,
			}
			for _, observer := range subscription.Observers {
				observer.Message <- notif
			}
		}
	}

	return nil
}

func (r repo) CreateCommentNotifFromConsumer(body []byte) error {
	var s *CommentEvent
	err := json.Unmarshal(body, &s)
	if err != nil {
		return err
	}
	
	m := []*model.Map{
		{
			Key: "commentBy",
			Value: s.Comment.CreatedBy,
		},
		{
			Key: "postID",
			Value: s.Comment.PostID,
		},
	}
	
	notifText := fmt.Sprintf("New Comment from %s", s.Comment.CreatedBy)
	notifEvent := NotifEvent{
		EventTime: time.Now(),
		EventType: "New Comment",
		Receiver:  s.Post.Username,
		Text:      notifText,
		Read:      false,
		Route:     "/profile",
		Params:    m,
	}

	insertedId, err := r.InsertNotifEvent(context.Background(), notifEvent)
	if err != nil {
		return err
	}

	subscription := r.Subscriptions[notifEvent.Receiver]
	if subscription != nil {
		notif := &model.Notif{
			ID:        insertedId,
			Event:     notifEvent.EventType,
			Timestamp: notifEvent.EventTime,
			Read:      notifEvent.Read,
			Receiver:  notifEvent.Receiver,
			Text:      notifEvent.Text,
			Params:    notifEvent.Params,
			Route:     notifEvent.Route,
		}
		for _, observer := range subscription.Observers {
			observer.Message <- notif
		}
	}

	return nil
}

func (r repo) CreateLikeNotifFromConsumer(body []byte) error {
	var s *PostEvent
	err := json.Unmarshal(body, &s)
	if err != nil {
		return err
	}

	m := []*model.Map{
		{
			Key: "postId",
			Value: s.PostID,
		},
	}

	notifText := fmt.Sprintf("You have a new like!")
	notifEvent := NotifEvent{
		EventTime: time.Now(),
		EventType: "New Like",
		Receiver:  s.Username,
		Text:      notifText,
		Read:      false,
		Route:     "/profile",
		Params:    m,
	}

	insertedId, err := r.InsertNotifEvent(context.Background(), notifEvent)
	if err != nil {
		return err
	}

	subscription := r.Subscriptions[notifEvent.Receiver]
	if subscription != nil {
		notif := &model.Notif{
			ID:        insertedId,
			Event:     notifEvent.EventType,
			Timestamp: notifEvent.EventTime,
			Read:      notifEvent.Read,
			Receiver:  notifEvent.Receiver,
			Text:      notifEvent.Text,
			Params:    notifEvent.Params,
			Route:     notifEvent.Route,
		}
		for _, observer := range subscription.Observers {
			observer.Message <- notif
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

func (r repo) GetNotifsByReceiver(ctx context.Context, receiver string) ([]*model.Notif, error) {

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
			Params:   notif.Params,
			Receiver: notif.Receiver,
			Read:     notif.Read,
			Route:    notif.Route,
			Text:     notif.Text,
			Timestamp: notif.EventTime,
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
		Subscriptions: map[string]*Message{},
	}, nil
}