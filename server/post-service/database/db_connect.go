package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const defaultMongoDBUrl = "mongodb://admin:admin@localhost:27017/"

/*
method for connecting to db
returns the mongodb client and an error, if one occurs
*/
func dbConnect() (*mongo.Client, error) {
	// get url from environment variable else take the default url
	url := os.Getenv("POST_MONGODB_URL")
	if url == "" {
		url = defaultMongoDBUrl
	}

	// create a new client with url
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	// create a context for connection timeout to mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to client
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// ping the client, so connection is really established
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}
