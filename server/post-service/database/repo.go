package database

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"posts-service/graph/model"
	"reflect"
	"strings"
)

type Repository interface {
	InsertPostEvent(post PostEvent) error
	InsertFile(base64File string) (*model.File, error)
	CreatePost(postEvent PostEvent, base64File string) (*model.Post, error)
	GetPosts(fetchLimit int) ([]*model.Post, error)
	GetCurrentPosts() []*model.Post
	GetPostById(id string) (int, *model.Post)
	RemovePost(postEvent PostEvent, index int) (string, error)
	EditPost(postEvent PostEvent) (string, error)
	LikePost(postEvent PostEvent) (string, error)
	UnlikePost(postEvent PostEvent) (string, error)
	AddComment(postEvent PostEvent) (string, error)
}

type repo struct {
	postCollection       *mongo.Collection
	fileCollection       *mongo.Collection
	bucket               *gridfs.Bucket
	currentPosts         []*model.Post
	lastFetchedEventTime string
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
		postCollection:       db.Collection("post-events"),
		fileCollection:       bucket.GetFilesCollection(),
		bucket:               bucket,
		lastFetchedEventTime: "",
		currentPosts:         make([]*model.Post, 0),
	}, nil
}

func (repo *repo) InsertPostEvent(postEvent PostEvent) error {
	_, err := repo.postCollection.InsertOne(ctx, postEvent)
	if err != nil {
		return err
	}

	return err
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

	err = repo.InsertPostEvent(postEvent)
	if err != nil {
		return nil, err
	}

	post := &model.Post{
		ID:          postEvent.PostID,
		Description: postEvent.Description,
		Data:        file,
		LikedBy:     postEvent.LikedBy,
		Comments:    postEvent.Comments,
	}

	// add to currentPosts
	repo.currentPosts = append([]*model.Post{post}, repo.currentPosts...)

	return post, nil
}

func (repo *repo) GetPosts(fetchLimit int) ([]*model.Post, error) {
	currentPosts := make([]*model.Post, 0)

	limit := int64(fetchLimit)

	// sort post-events by descending event-time (the newest first) and set fetch limit
	opts := options.Find()
	opts.SetSort(bson.D{{"event_time", -1}})
	opts.Limit = &limit

	// check if it is the first time getting data
	key := "$lt"
	if repo.lastFetchedEventTime == "" {
		key = "$gt"
	}

	// get all post events with event_type = "CreatePost" sorted by event_time
	cursor, err := repo.postCollection.Find(ctx, bson.D{
		{"event_type", "CreatePost"},
		{"event_time", bson.D{{key, repo.lastFetchedEventTime}}},
	}, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var postEvent PostEvent
		if err = cursor.Decode(&postEvent); err != nil {
			return nil, err
		}

		// get file content for post
		var buf bytes.Buffer
		_, err := repo.bucket.DownloadToStreamByName(postEvent.FileID, &buf)
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

		// add new post to output for getPosts
		currentPosts = append(currentPosts, &model.Post{
			ID:          postEvent.PostID,
			Description: postEvent.Description,
			Data: &model.File{
				Name:        postEvent.FileID,
				Content:     string(buf.Bytes()),
				ContentType: fmt.Sprint(contentType.Interface()),
			},
			LikedBy:  postEvent.LikedBy,
			Comments: postEvent.Comments,
		})

		// update last fetched event time
		repo.lastFetchedEventTime = postEvent.EventTime
	}

	// TODO get edited data

	// update runtime data
	repo.currentPosts = append(repo.currentPosts, currentPosts...)

	return currentPosts, nil
}

func (repo *repo) GetCurrentPosts() []*model.Post {
	return repo.currentPosts
}

func (repo *repo) GetPostById(id string) (int, *model.Post) {
	for i, post := range repo.currentPosts {
		if post.ID == id {
			return i, post
		}
	}

	return -1, nil
}

func (repo *repo) RemovePost(postEvent PostEvent, index int) (string, error) {
	/*// remove from currentPosts
	repo.currentPosts = append(repo.currentPosts[:index], repo.currentPosts[index+1:]...)

	// delete all events relating to the id
	sqlQuery := `delete from "post-events" where "postId" = $1`

	_, err := repo.client.Exec(sqlQuery, postEvent.PostID)
	if err != nil {
		return "failed", err
	}

	// new current post events
	err = repo.InsertPostEvent(postEvent)
	if err != nil {
		return "failed", err
	}*/

	return "success", nil
}

func (repo *repo) EditPost(postEvent PostEvent) (string, error) {
	err := repo.InsertPostEvent(postEvent)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *repo) LikePost(postEvent PostEvent) (string, error) {
	err := repo.InsertPostEvent(postEvent)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *repo) UnlikePost(postEvent PostEvent) (string, error) {
	err := repo.InsertPostEvent(postEvent)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *repo) AddComment(postEvent PostEvent) (string, error) {
	err := repo.InsertPostEvent(postEvent)
	if err != nil {
		return "failure", nil
	}

	return "success", nil
}
