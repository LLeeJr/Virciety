package database

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
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
	AddProfilePicture(ctx context.Context, profilePictureEvent ProfilePictureEvent, input model.AddProfilePicture) (string, error)
	InsertFile(base64 string) (*model.File, error)
	GetProfilePicture(ctx context.Context, fileID *string) (string, error)
	RemoveProfilePicture(ctx context.Context, profilePictureEvent ProfilePictureEvent) (string, error)
}

type repo struct {
	userCollection           *mongo.Collection
	fileCollection           *mongo.Collection
	chunksCollection         *mongo.Collection
	profilePictureCollection *mongo.Collection
	bucket                   *gridfs.Bucket
}

func (r repo) RemoveProfilePicture(ctx context.Context, profilePictureEvent ProfilePictureEvent) (string, error) {

	type file struct {
		ID         primitive.ObjectID `bson:"_id"`
		FileName   string             `bson:"filename"`
	}

	var result *file
	if err := r.fileCollection.FindOne(ctx, bson.D{
		{"filename", profilePictureEvent.FileId},
	}).Decode(&result); err != nil {
		return "error during finding profile picture entry", err
	}

	_, err := r.chunksCollection.DeleteMany(ctx, bson.D{
		{"files_id", result.ID},
	})
	if err != nil {
		return "error during deleting file chunks", err
	}

	_, err = r.fileCollection.DeleteOne(ctx, bson.D{
		{"filename", profilePictureEvent.FileId},
	})
	if err != nil {
		return "error during deleting profile picture", err
	}

	_, err = r.profilePictureCollection.DeleteOne(ctx, bson.D{
		{"fileid", profilePictureEvent.FileId},
	})
	if err != nil {
		return "error during deleting profile picture event", err
	}

	user, err := r.GetUserByName(ctx, &profilePictureEvent.Username)
	if err != nil {
		return "error while retrieving user profile", err
	}

	objectID, err := primitive.ObjectIDFromHex(user.ID)
	query := bson.M{
		"_id": objectID,
	}

	update := bson.D{
		{"$set", bson.D{{"fileId", ""}}},
	}

	_, err = r.userCollection.UpdateOne(ctx,
		query,
		update,
	)

	return "removed profile picture successfully", err
}

func (r repo) GetProfilePicture(ctx context.Context, fileID *string) (string, error) {
	var buf bytes.Buffer
	_, err := r.bucket.DownloadToStreamByName(*fileID, &buf)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func (r repo) InsertFile(base64 string) (*model.File, error) {
	fileName := uuid.NewString()

	properties := strings.Split(base64, ";base64,")
	contentType := strings.Split(properties[0], ":")

	uploadOpts := options.GridFSUpload().
		SetMetadata(bson.D{{"contentType", contentType[1]}})

	uploadStream, err := r.bucket.OpenUploadStream(fileName, uploadOpts)
	if err != nil {
		return nil, err
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write([]byte(base64))
	if err != nil {
		return nil, err
	}

	return &model.File{
		Name:        fileName,
		Content:     base64,
		ContentType: contentType[1],
	}, nil
}


func (r repo) AddProfilePicture(ctx context.Context, profilePictureEvent ProfilePictureEvent, input model.AddProfilePicture) (string, error) {
	user, err := r.GetUserByName(ctx, &input.Username)
	if err != nil {
		return "error while retrieving user profile", err
	}

	if user.ProfilePictureID != "" {
		// profile picture already exists in db, which is why it needs to be deleted at first
		profilePictureEvent.FileId = user.ProfilePictureID
		_, err = r.RemoveProfilePicture(ctx, profilePictureEvent)
		if err != nil {
			return "error during removing old profile picture", err
		}
	}

	file, err := r.InsertFile(input.Data)
	if err != nil {
		return "error while inserting file in db", err
	}

	profilePictureEvent.FileId = file.Name

	_, err = r.InsertProfilePictureEvent(ctx, profilePictureEvent)
	if err != nil {
		return "error while inserting file related data in db", err
	}

	objectID, err := primitive.ObjectIDFromHex(user.ID)
	query := bson.M{
		"_id": objectID,
	}

	update := bson.D{
		{"$set", bson.D{{"fileId", file.Name}}},
	}

	_, err = r.userCollection.UpdateOne(ctx,
		query,
		update,
	)
	if err != nil {
		return "error during updating user profile", err
	}

	return file.Name, nil

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
	FileId    string             `bson:"fileId"`
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
		ProfilePictureID: result.FileId,
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
			ProfilePictureID: userResult.FileId,
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
		ProfilePictureID: result.FileId,
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
		ProfilePictureID: "",
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

func (r repo) InsertProfilePictureEvent(ctx context.Context, profilePictureEvent ProfilePictureEvent) (string, error) {
	inserted, err := r.profilePictureCollection.InsertOne(ctx, profilePictureEvent)
	if err != nil {
		return "error while inserting profile picture related data into db", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), nil
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
		userCollection:   db.Collection("user-events"),
		profilePictureCollection: db.Collection("profile-picture-events"),
		fileCollection:   bucket.GetFilesCollection(),
		chunksCollection: bucket.GetChunksCollection(),
		bucket:           bucket,
	}, nil
}