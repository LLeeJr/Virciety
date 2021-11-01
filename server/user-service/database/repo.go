package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type Repository interface {

}

type repo struct {
	userCollection *mongo.Collection
	bucket *gridfs.Bucket
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