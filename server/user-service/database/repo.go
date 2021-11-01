package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"user-service/graph/model"
)

type Repository interface {
	CreateUser(ctx context.Context, userEvent UserEvent) (*model.User, error)
	InsertUserEvent(ctx context.Context, userEvent UserEvent) (string, error)
}

type repo struct {
	userCollection *mongo.Collection
	bucket *gridfs.Bucket
}

func (r repo) CreateUser(ctx context.Context, userEvent UserEvent) (*model.User, error) {
	_, err := r.userCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	
	insertedId, err := r.InsertUserEvent(ctx, userEvent)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:        insertedId,
		Username:  userEvent.Username,
		FirstName: userEvent.FirstName,
		LastName:  userEvent.LastName,
		Follows:   userEvent.Follows,
	}

	return user, nil
}

func (r repo) InsertUserEvent(ctx context.Context, userEvent UserEvent) (string, error) {
	inserted, err := r.userCollection.InsertOne(ctx, userEvent)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), err
}

func NewRepo() (Repository, error) {
	client, err := Connect()
	if err != nil {
		return nil, err
	}

	db := client.Database("user-service")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, err
	}

	return &repo{
		userCollection: db.Collection("user-events"),
		bucket:         bucket,
	}, nil
}