package database

import (
	"context"
	"dm-service/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type UniqueUser struct {
	UserName string `json:"userName"`
}

type Repository interface {
	GetRoom(ctx context.Context, name string) (*model.Chatroom, error)
	GetRoomsByUser(ctx context.Context, userName string) ([]*model.Chatroom, error)
	CreateDm(ctx context.Context, dmEvent DmEvent) (*model.Dm, error)
	CreateRoom(ctx context.Context, roomEvent ChatroomEvent) (*model.Chatroom, error)
	UpdateRoom(ctx context.Context, room *model.Chatroom) (string, error)
	InsertDmEvent(ctx context.Context, dmEvent DmEvent) (string, error)
	InsertRoomEvent(ctx context.Context, roomEvent ChatroomEvent) (string, error)
}

type repo struct {
	dmCollection   *mongo.Collection
	roomCollection *mongo.Collection
	bucket         *gridfs.Bucket
}

func (r repo) UpdateRoom(ctx context.Context, room *model.Chatroom) (string, error) {
	objID, err := primitive.ObjectIDFromHex(room.ID)
	if err != nil {
		return "", err
	}

	query := bson.M{
		"_id": objID,
	}

	update := bson.M{
		"$set": bson.M{
			"member": room.Member,
		},
	}

	_, err = r.roomCollection.UpdateOne(ctx,
		query,
		update,
	)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func (r repo) CreateRoom(ctx context.Context, roomEvent ChatroomEvent) (*model.Chatroom, error) {
	insertedId, err := r.InsertRoomEvent(ctx, roomEvent)
	if err != nil {
		return nil, err
	}

	room := &model.Chatroom{
		ID:       insertedId,
		Member:   roomEvent.Member,
		Name:     roomEvent.Name,
	}

	return room, nil
}

func (r repo) InsertRoomEvent(ctx context.Context, room ChatroomEvent) (string, error) {
	inserted, err := r.roomCollection.InsertOne(ctx, room)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), err
}

func (r repo) InsertDmEvent(ctx context.Context, dmEvent DmEvent) (string, error) {
	inserted, err := r.dmCollection.InsertOne(ctx, dmEvent)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), err
}

type Chatroom struct {
	ID        primitive.ObjectID `bson:"_id"`
	Member    []string           `bson:"member"`
	Name      string             `bson:"name"`
	EventType string             `bson:"eventtype"`
}

func (r repo) GetRoom(ctx context.Context, name string) (*model.Chatroom, error) {

	var result *Chatroom
	if err := r.roomCollection.FindOne(ctx, bson.D{
		{"name", name},
	}).Decode(&result); err != nil {
		return nil, err
	}

	room := &model.Chatroom{
		ID:     result.ID.Hex(),
		Member: result.Member,
		Name:   result.Name,
	}

	return room, nil
}

func (r repo) GetRoomsByUser(ctx context.Context, userName string) ([]*model.Chatroom, error) {

	var result []*Chatroom
	cursor, err := r.roomCollection.Find(
		ctx,
		bson.D{
			{"member", userName},
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	rooms := make([]*model.Chatroom, 0)
	for _, chatroom := range result {
		rooms = append(rooms, &model.Chatroom{
			ID:     chatroom.ID.Hex(),
			Member: chatroom.Member,
			Name:   chatroom.Name,
		})
	}

	return rooms, err
}

func (r repo) CreateDm(ctx context.Context, dmEvent DmEvent) (*model.Dm, error) {
	insertedId, err := r.InsertDmEvent(ctx, dmEvent)
	if err != nil {
		return nil, err
	}

	dm := &model.Dm{
		ChatroomID: dmEvent.ChatroomId,
		CreatedAt:  dmEvent.CreatedAt,
		CreatedBy:  dmEvent.CreatedBy,
		ID:         insertedId,
		Msg:        dmEvent.Msg,
	}

	return dm, err
}

func NewRepo() (Repository, error) {
	client, err := Connect()
	if err != nil {
		return nil, err
	}

	db := client.Database("dm-service")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, err
	}

	return &repo{
		dmCollection: db.Collection("dm-events"),
		roomCollection: db.Collection("room-events"),
		bucket: bucket,
	}, nil
}


