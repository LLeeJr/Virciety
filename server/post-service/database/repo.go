package database

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"post-service/graph/model"
	"reflect"
	"strings"
)

type Repository interface {
	InsertPostEvent(ctx context.Context, post PostEvent) (string, error)
	InsertFile(ctx context.Context, base64File string) (*model.File, error)
	CreatePost(ctx context.Context, postEvent PostEvent, base64File string) (*model.Post, error)
	GetPosts(ctx context.Context, id string, fetchLimit int, filter *string) ([]*model.Post, error)
	RemovePost(ctx context.Context, postEvent PostEvent) (string, error)
	EditPost(ctx context.Context, postEvent PostEvent) (string, error)
	LikePost(ctx context.Context, postEvent PostEvent) error
	GetData(ctx context.Context, fileID string) (string, error)
	GetPost(ctx context.Context, id string) (*model.Post, error)
}

type Repo struct {
	postCollection   *mongo.Collection
	fileCollection   *mongo.Collection
	chunksCollection *mongo.Collection
	bucket           *gridfs.Bucket
}

func NewRepo() (Repository, error) {
	client, err := dbConnect()
	if err != nil {
		return nil, err
	}

	db := client.Database("post-service")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, err
	}

	return &Repo{
		postCollection:   db.Collection("post-events"),
		fileCollection:   bucket.GetFilesCollection(),
		chunksCollection: bucket.GetChunksCollection(),
		bucket:           bucket,
	}, nil
}

func (repo *Repo) InsertPostEvent(ctx context.Context, postEvent PostEvent) (string, error) {
	inserted, err := repo.postCollection.InsertOne(ctx, postEvent)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), err
}

func (repo *Repo) InsertFile(_ context.Context, base64File string) (*model.File, error) {
	fileName := uuid.NewString()

	properties := strings.Split(base64File, ";base64,")
	contentType := strings.Split(properties[0], ":")

	uploadOpts := options.GridFSUpload().
		SetMetadata(bson.D{{"contentType", contentType[1]}})

	uploadStream, err := repo.bucket.OpenUploadStream(fileName, uploadOpts)
	if err != nil {
		return nil, err
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write([]byte(base64File))
	if err != nil {
		return nil, err
	}

	return &model.File{
		Name:        fileName,
		Content:     base64File,
		ContentType: contentType[1],
	}, nil
}

func (repo *Repo) CreatePost(ctx context.Context, postEvent PostEvent, base64File string) (*model.Post, error) {
	file, err := repo.InsertFile(ctx, base64File)
	if err != nil {
		return nil, err
	}

	postEvent.FileID = file.Name

	insertedID, err := repo.InsertPostEvent(ctx, postEvent)
	if err != nil {
		return nil, err
	}

	post := &model.Post{
		ID:          insertedID,
		Description: postEvent.Description,
		Data:        file,
		Username:    postEvent.Username,
		LikedBy:     postEvent.LikedBy,
		Comments:    postEvent.Comments,
	}

	return post, nil
}

func (repo *Repo) GetPost(ctx context.Context, id string) (*model.Post, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var postEvent *PostEvent
	err = repo.postCollection.FindOne(ctx, bson.D{
		{"_id", objID},
	}).Decode(&postEvent)
	if err != nil {
		return nil, err
	}

	// get file contentType
	var file bson.M
	if err = repo.fileCollection.FindOne(ctx, bson.M{"filename": postEvent.FileID}).Decode(&file); err != nil {
		return nil, err
	}

	// convert interface{} to map and get contentType
	var contentType reflect.Value
	metaData := file["metadata"]
	v := reflect.ValueOf(metaData)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			contentType = v.MapIndex(key)
		}
	}

	post := &model.Post{
		ID:          postEvent.ID.Hex(),
		Description: postEvent.Description,
		Data: &model.File{
			Name:        postEvent.FileID,
			ContentType: fmt.Sprint(contentType.Interface()),
		},
		Username: postEvent.Username,
		LikedBy:  postEvent.LikedBy,
		Comments: postEvent.Comments,
	}

	max := int64(1)
	// Sort event_time and get one element which will be the most recent edited post in relation to liked, unliked and description
	opts := options.Find()
	opts.SetSort(bson.D{{"event_time", -1}})
	opts.Limit = &max
	cursor, err := repo.postCollection.Find(ctx, bson.D{
		{"id", post.ID},
		{"event_type", bson.D{
			{"$in", bson.A{"EditPost", "LikePost", "UnlikePost"}},
		}},
	}, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var postEvent PostEvent
		if err = cursor.Decode(&postEvent); err != nil {
			return nil, err
		}

		// Add editable data
		post.LikedBy = postEvent.LikedBy
		post.Description = postEvent.Description
		post.Comments = postEvent.Comments
	}

	return post, nil
}

func (repo *Repo) GetPosts(ctx context.Context, id string, fetchLimit int, filter *string) ([]*model.Post, error) {
	currentPosts := make([]*model.Post, 0)
	limit := int64(fetchLimit)

	lastFetchedEventTime := ""
	if id != "" {
		projection := bson.D{
			{"_id", 0},
			{"event_time", 1},
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}

		var result bson.M
		err = repo.postCollection.FindOne(ctx, bson.D{
			{"_id", objID},
			{"event_type", "CreatePost"},
		}, options.FindOne().SetProjection(projection)).Decode(&result)
		if err != nil {
			return nil, err
		}

		lastFetchedEventTime = fmt.Sprint(result["event_time"])
	}

	// sort post-events by descending event-time (the newest first) and set fetch limit
	opts := options.Find()
	opts.SetSort(bson.D{{"event_time", -1}})
	opts.Limit = &limit

	// check if it is the first time getting data
	key := "$lt"
	if lastFetchedEventTime == "" {
		key = "$gt"
	}

	find := bson.D{
		{"event_type", "CreatePost"},
		{"event_time", bson.D{{key, lastFetchedEventTime}}}}

	if filter != nil {
		find = bson.D{
			{"event_type", "CreatePost"},
			{"event_time", bson.D{{key, lastFetchedEventTime}}},
			{"username", filter}}
	}

	// get all post events with event_type = "CreatePost" sorted by event_time
	cursor, err := repo.postCollection.Find(ctx, find, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var postEvent PostEvent
		if err = cursor.Decode(&postEvent); err != nil {
			return nil, err
		}

		// get file contentType
		var file bson.M
		if err = repo.fileCollection.FindOne(ctx, bson.M{"filename": postEvent.FileID}).Decode(&file); err != nil {
			return nil, err
		}

		// convert interface{} to map and get contentType
		var contentType reflect.Value
		metaData := file["metadata"]
		v := reflect.ValueOf(metaData)
		if v.Kind() == reflect.Map {
			for _, key := range v.MapKeys() {
				contentType = v.MapIndex(key)
			}
		}

		// add new post to output for getPosts
		currentPosts = append(currentPosts, &model.Post{
			ID:          postEvent.ID.Hex(),
			Description: postEvent.Description,
			Data: &model.File{
				Name:        postEvent.FileID,
				ContentType: fmt.Sprint(contentType.Interface()),
			},
			Username: postEvent.Username,
			LikedBy:  postEvent.LikedBy,
			Comments: postEvent.Comments,
		})
	}

	max := int64(1)
	for _, post := range currentPosts {
		// Sort event_time and get one element which will be the most recent edited post in relation to liked, unliked and description
		opts.Limit = &max
		cursor, err = repo.postCollection.Find(ctx, bson.D{
			{"id", post.ID},
			{"event_type", bson.D{
				{"$in", bson.A{"EditPost", "LikePost", "UnlikePost"}},
			}},
		}, opts)
		if err != nil {
			return nil, err
		}

		for cursor.Next(ctx) {
			var postEvent PostEvent
			if err = cursor.Decode(&postEvent); err != nil {
				return nil, err
			}

			// Add editable data
			post.LikedBy = postEvent.LikedBy
			post.Description = postEvent.Description
			post.Comments = postEvent.Comments
		}
	}

	return currentPosts, nil
}

func (repo *Repo) GetData(_ context.Context, fileID string) (string, error) {
	// get file content for post
	var buf bytes.Buffer
	_, err := repo.bucket.DownloadToStreamByName(fileID, &buf)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func (repo *Repo) RemovePost(ctx context.Context, postEvent PostEvent) (string, error) {
	// get fileID for deleting all chunks
	projection := bson.D{
		{"_id", 1},
	}

	var result bson.M
	if err := repo.fileCollection.FindOne(ctx, bson.M{"filename": postEvent.FileID}, options.FindOne().SetProjection(projection)).Decode(&result); err != nil {
		return "failed", err
	}

	fileID := result["_id"].(primitive.ObjectID)

	// delete file ref from fileCollection
	_, err := repo.fileCollection.DeleteOne(ctx, bson.D{
		{"_id", fileID},
	})
	if err != nil {
		return "failed", err
	}

	// delete allChunks from chunksCollection
	_, err = repo.chunksCollection.DeleteMany(ctx, bson.D{
		{"files_id", fileID},
	})
	if err != nil {
		return "failed", err
	}

	// convert hex-string into primitive.objectID
	objID, err := primitive.ObjectIDFromHex(postEvent.PostID)
	if err != nil {
		return "failed", err
	}

	// delete that one CreatePost-Event
	_, err = repo.postCollection.DeleteOne(ctx, bson.D{
		{"_id", objID},
		{"event_type", "CreatePost"},
	})
	if err != nil {
		return "failed", err
	}

	// delete all other events
	_, err = repo.postCollection.DeleteMany(ctx, bson.D{
		{"id", postEvent.PostID},
		{"event_type", bson.D{
			{"$in", bson.A{"EditPost", "LikePost", "UnlikePost"}},
		}},
	})
	if err != nil {
		return "failed", err
	}

	// new current post events
	_, err = repo.InsertPostEvent(ctx, postEvent)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *Repo) EditPost(ctx context.Context, postEvent PostEvent) (string, error) {
	_, err := repo.InsertPostEvent(ctx, postEvent)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *Repo) LikePost(ctx context.Context, postEvent PostEvent) error {
	_, err := repo.InsertPostEvent(ctx, postEvent)
	if err != nil {
		return err
	}

	return nil
}
