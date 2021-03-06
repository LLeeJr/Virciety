package database

import (
	"context"
	"dm-service/graph/model"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type UniqueUser struct {
	UserName string `json:"userName"`
}

type Repository interface {
	GetRoom(ctx context.Context, roomName string, id string) (*model.Chatroom, error)
	GetRoomsByUser(ctx context.Context, userName string) ([]*model.Chatroom, error)
	CreateDm(ctx context.Context, dmEvent DmEvent) (*model.Dm, error)
	CreateRoom(ctx context.Context, roomEvent ChatroomEvent) (*model.Chatroom, error)
	UpdateRoom(ctx context.Context, room *model.Chatroom) (string, error)
	InsertDmEvent(ctx context.Context, dmEvent DmEvent) (string, error)
	InsertRoomEvent(ctx context.Context, roomEvent ChatroomEvent) (string, error)
	GetMessagesFromRoom(ctx context.Context, id string) ([]*model.Dm, error)
	DeleteRoom(ctx context.Context, id string) (string, error)
	GetDirectRoom(ctx context.Context, user1 string, user2 string) (*model.Chatroom, error)
}

type repo struct {
	dmCollection   *mongo.Collection
	roomCollection *mongo.Collection
	bucket         *gridfs.Bucket
}

type Chatroom struct {
	EventType string             `bson:"eventtype"`
	ID        primitive.ObjectID `bson:"_id"`
	Member    []string           `bson:"member"`
	Name      string             `bson:"name"`
	Owner     string             `bson:"owner"`
	IsDirect  bool               `bson:"isdirect"`
}

// GetDirectRoom retrieves a direct-chat chatroom for two given userNames from the database
func (r repo) GetDirectRoom(ctx context.Context, user1 string, user2 string) (*model.Chatroom, error) {
	order1 := []string{user1, user2}
	order2 := []string{user2, user1}
	query := bson.M{
		"$or": []interface{}{
			bson.M{
				"$and": []interface{}{
					bson.M{"membersize": bson.M{"$eq": 2}},
					bson.M{"member": order1},
					bson.M{"isdirect": true},
				},
			},
			bson.M{
				"$and": []interface{}{
					bson.M{"membersize": bson.M{"$eq": 2}},
					bson.M{"member": order2},
					bson.M{"isdirect": true},
				},
			},
		},
	}

	var rooms []*Chatroom
	cursor, err := r.roomCollection.Find(ctx, query)

	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &rooms)
	if err != nil {
		return nil, err
	}

	if len(rooms) != 0 {
		chatroom := &model.Chatroom{
			ID:     rooms[0].ID.Hex(),
			IsDirect: rooms[0].IsDirect,
			Member: rooms[0].Member,
			Name:   rooms[0].Name,
			Owner:  rooms[0].Owner,
		}

		return chatroom, nil
	}

	return nil, errors.New("no room found")
}

// DeleteRoom removes a chatroom from the database for a given id
func (r repo) DeleteRoom(ctx context.Context, id string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(id)

	result, err := r.roomCollection.DeleteOne(ctx, bson.D{
		{"_id", objID},
	})
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("Deleted %v room!", result.DeletedCount)

	_, err = r.dmCollection.DeleteMany(ctx, bson.D{
		{"chatroomid", id},
	})

	return msg, err
}

// GetMessagesFromRoom retrieves all dms from the database for a given roomID
func (r repo) GetMessagesFromRoom(ctx context.Context, id string) ([]*model.Dm, error) {
	var result []*DmEvent
	cursor, err := r.dmCollection.Find(
		ctx,
		bson.D{
			{"chatroomid", id},
		},
	)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	dms := make([]*model.Dm, 0)
	for _, event := range result {
		dms = append(dms, &model.Dm{
			ChatroomID: id,
			CreatedAt:  event.CreatedAt,
			CreatedBy:  event.CreatedBy,
			Msg:        event.Msg,
		})
	}

	return dms, err
}

// UpdateRoom updates an existing chatroom inside the database
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
			"member":     room.Member,
			"membersize": len(room.Member),
			"owner":      room.Owner,
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

// CreateRoom creates a new chatroom in the database (stored as event)
func (r repo) CreateRoom(ctx context.Context, roomEvent ChatroomEvent) (*model.Chatroom, error) {
	insertedId, err := r.InsertRoomEvent(ctx, roomEvent)
	if err != nil {
		return nil, err
	}

	modelRoom := &model.Chatroom{
		ID:       insertedId,
		IsDirect: *roomEvent.IsDirect,
		Member:   roomEvent.Member,
		Name:     roomEvent.Name,
		Owner:    roomEvent.Owner,
	}

	return modelRoom, nil
}

// InsertRoomEvent is a helper-function for inserting a ChatroomEvent in the database
func (r repo) InsertRoomEvent(ctx context.Context, room ChatroomEvent) (string, error) {
	inserted, err := r.roomCollection.InsertOne(ctx, room)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), err
}

// InsertDmEvent is a helper-function for inserting a DmEvent in the database
func (r repo) InsertDmEvent(ctx context.Context, dmEvent DmEvent) (string, error) {
	inserted, err := r.dmCollection.InsertOne(ctx, dmEvent)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), err
}

// GetRoom retrieves a chatroom from the database for a given roomName and a respective roomID
func (r repo) GetRoom(ctx context.Context, roomName string, id string) (*model.Chatroom, error) {
	var objID primitive.ObjectID
	if id != "" {
		var err error
		objID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
	}

	query := bson.M{
		"$or": []interface{}{
			bson.M{"_id": objID},
			bson.M{"name": roomName},
		},
	}

	var result *Chatroom
	if err := r.roomCollection.FindOne(ctx, query).Decode(&result); err != nil {
		return nil, err
	}

	return &model.Chatroom{
		ID:     result.ID.Hex(),
		IsDirect: result.IsDirect,
		Member: result.Member,
		Name:   result.Name,
		Owner:  result.Owner,
	}, nil
}

// GetRoomsByUser retrieves all chatroom-elements from the database for a given userName
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
			Owner:  chatroom.Owner,
			IsDirect: chatroom.IsDirect,
		})
	}

	return rooms, err
}

// CreateDm creates a new dm in the database (stored as event)
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

// NewRepo creates a new Repository instance for the given database
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


