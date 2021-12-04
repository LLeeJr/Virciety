package database

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"user-service/graph/model"
)

type Repository interface {
	CreateUser(ctx context.Context, userEvent UserEvent) (*model.User, error)
	InsertUserEvent(ctx context.Context, userEvent UserEvent) (string, error)
	GetUserByID(ctx context.Context, id *string) (*model.User, error)
	GetUserByName(ctx context.Context, name *string) (*model.User, error)
	AddFollow(ctx context.Context, id *string, username *string, add *string) (*model.User, error)
	RemoveFollow(ctx context.Context, id *string, username *string, remove *string) (*model.User, error)
	FindUsersWithName(ctx context.Context, name *string) ([]*model.User, error)
}

type repo struct {
	userCollection *mongo.Collection
	bucket *gridfs.Bucket
}

func (r repo) RemoveFollow(ctx context.Context, id *string, username *string, remove *string) (*model.User, error) {
	if *username == *remove {
		return nil, errors.New("user can not unfollow himself")
	}

	objectID, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"_id": objectID,
	}

	update := bson.M{
		"$pull": bson.M{
			"follows": remove,
		},
	}

	updated, err := r.userCollection.UpdateOne(ctx,
		query,
		update,
	)

	if err != nil {
		return nil, err
	}

	if updated.ModifiedCount == 0 {
		return nil, errors.New("could not update user")
	}

	query = bson.M{
		"username": remove,
	}

	update = bson.M{
		"$pull": bson.M{
			"followers": username,
		},
	}

	updated, err = r.userCollection.UpdateOne(ctx,
		query,
		update,
	)

	if err != nil {
		return nil, err
	}

	if updated.ModifiedCount == 0 {
		return nil, errors.New("could not update user")
	}

	user, err := r.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r repo) AddFollow(ctx context.Context, id *string, username *string, add *string) (*model.User, error) {
	if *username == *add {
		return nil, errors.New("user can not follow himself")
	}

	objectID, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"_id": objectID,
	}

	update := bson.M{
		"$addToSet": bson.M{
			"follows": add,
		},
	}

	updated, err := r.userCollection.UpdateOne(ctx,
		query,
		update,
	)
	if err != nil {
		return nil, err
	}

	if updated.ModifiedCount == 0 {
		return nil, errors.New("could not update user")
	}

	query = bson.M{
		"username": add,
	}

	update = bson.M{
		"$addToSet": bson.M{
			"followers": username,
		},
	}

	updated, err = r.userCollection.UpdateOne(ctx,
		query,
		update,
	)

	if err != nil {
		return nil, err
	}

	if updated.ModifiedCount == 0 {
		return nil, errors.New("could not update user")
	}

	user, err := r.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type UserByName struct {
	EventType string             `bson:"eventtype"`
	EventTime time.Time          `bson:"eventtime"`
	FirstName string             `bson:"firstname"`
	Follows   []string           `bson:"follows"`
	Followers []string           `bson:"followers"`
	ID        primitive.ObjectID `bson:"_id"`
	LastName  string             `bson:"lastname"`
	Username  string             `bson:"username"`
}

func (r repo) GetUserByName(ctx context.Context, name *string) (*model.User, error) {

	var result *UserByName
	if err := r.userCollection.FindOne(ctx, bson.D{
		{"username", name},
	}).Decode(&result); err != nil {
		return nil, err
	}

	user := &model.User{
		ID:        result.ID.Hex(),
		Username:  result.Username,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Follows:   result.Follows,
		Followers: result.Followers,
	}

	return user, nil
}

func (r repo) FindUsersWithName(ctx context.Context, name *string) ([]*model.User, error) {

	findOptions := options.Find().SetLimit(10)

	pattern := fmt.Sprint("^", *name)
	regEx := primitive.Regex{Pattern: pattern, Options: "i"}

	filter := bson.M{"username": bson.M{"$regex": regEx}}

	cursor, err := r.userCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var result []*UserByName
	users := make([]*model.User, 0)
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	for _, userResult := range result {
		user := &model.User{
			ID:        userResult.ID.Hex(),
			Username:  userResult.Username,
			FirstName: userResult.FirstName,
			LastName:  userResult.LastName,
			Follows:   userResult.Follows,
			Followers: userResult.Followers,
		}

		users = append(users, user)
	}

	return users, nil
}

func (r repo) GetUserByID(ctx context.Context, id *string) (*model.User, error) {

	var result *UserEvent
	convertedId, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return nil, err
	}
	if err := r.userCollection.FindOne(ctx, bson.D{
		{"_id", convertedId},
	}).Decode(&result); err != nil {
		return nil, err
	}

	user := &model.User{
		ID:        *id,
		Username:  result.Username,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Follows:   result.Follows,
		Followers: result.Followers,
	}

	return user, nil
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
		Followers: userEvent.Followers,
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