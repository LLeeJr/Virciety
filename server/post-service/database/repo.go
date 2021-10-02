package database

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"posts-service/graph/model"
	"reflect"
	"strings"
)

type Repository interface {
	InsertPostEvent(post PostEvent) (string, error)
	InsertFile(base64File string) (*model.File, error)
	CreatePost(postEvent PostEvent, base64File string) (*model.Post, error)
	GetPosts(id string, fetchLimit int) ([]*model.Post, error)
	RemovePost(postEvent PostEvent) (string, error)
	EditPost(postEvent PostEvent) (string, error)
	LikePost(postEvent PostEvent) error
	AddComment(postEvent PostEvent) (string, error)
	GetData(fileID string) (string, error)
}

type repo struct {
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

	return &repo{
		postCollection:   db.Collection("post-events"),
		fileCollection:   bucket.GetFilesCollection(),
		chunksCollection: bucket.GetChunksCollection(),
		bucket:           bucket,
	}, nil
}

func (repo *repo) InsertPostEvent(postEvent PostEvent) (string, error) {
	inserted, err := repo.postCollection.InsertOne(ctx, postEvent)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), err
}

func (repo *repo) InsertFile(base64File string) (*model.File, error) {
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

func (repo *repo) CreatePost(postEvent PostEvent, base64File string) (*model.Post, error) {
	file, err := repo.InsertFile(base64File)
	if err != nil {
		return nil, err
	}

	postEvent.FileID = file.Name

	insertedID, err := repo.InsertPostEvent(postEvent)
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

func (repo *repo) GetPosts(id string, fetchLimit int) ([]*model.Post, error) {
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

	// get all post events with event_type = "CreatePost" sorted by event_time
	cursor, err := repo.postCollection.Find(ctx, bson.D{
		{"event_type", "CreatePost"},
		{"event_time", bson.D{{key, lastFetchedEventTime}}},
	}, opts)
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

func (repo *repo) GetData(fileID string) (string, error) {
	// get file content for post
	var buf bytes.Buffer
	_, err := repo.bucket.DownloadToStreamByName(fileID, &buf)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func (repo *repo) RemovePost(postEvent PostEvent) (string, error) {
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
	_, err = repo.InsertPostEvent(postEvent)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *repo) EditPost(postEvent PostEvent) (string, error) {
	_, err := repo.InsertPostEvent(postEvent)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *repo) LikePost(postEvent PostEvent) error {
	_, err := repo.InsertPostEvent(postEvent)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repo) AddComment(postEvent PostEvent) (string, error) {
	_, err := repo.InsertPostEvent(postEvent)
	if err != nil {
		return "failure", nil
	}

	return "success", nil
}
