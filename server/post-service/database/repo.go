package database

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"posts-service/graph/model"
	"strings"
)

type Repository interface {
	InsertPostEvent(post PostEvent) error
	InsertFile(base64File string) (*model.File, error)
	CreatePost(postEvent PostEvent, base64File string) (*model.Post, error)
	GetPosts() ([]*model.Post, error)
	GetCurrentPosts() []*model.Post
	GetPostById(id string) (int, *model.Post)
	RemovePost(postEvent PostEvent, index int) (string, error)
	EditPost(postEvent PostEvent) (string, error)
	LikePost(postEvent PostEvent) (string, error)
	UnlikePost(postEvent PostEvent) (string, error)
	AddComment(postEvent PostEvent) (string, error)
}

type repo struct {
	postCollection *mongo.Collection
	bucket         *gridfs.Bucket
	currentPosts   []*model.Post
	currentEventId int
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
		postCollection: db.Collection("post-events"),
		bucket:         bucket,
		currentEventId: 0,
		currentPosts:   make([]*model.Post, 0),
	}, nil
}

func (repo *repo) InsertPostEvent(postEvent PostEvent) (err error) {
	/*sqlQuery := `INSERT INTO "post-events" ("postId", "eventTime", "eventType", username, description, data, liked, comments)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	newTime, _ := time.Parse("2006-01-02 15:04:05", postEvent.EventTime)

	id := 0

	data := postEvent.Data.ID + "." + postEvent.Data.ContentType

	err = repo.DB.QueryRow(sqlQuery, postEvent.PostID, newTime, postEvent.EventType, postEvent.Username, postEvent.Description, data,
		pq.Array(postEvent.LikedBy), pq.Array(postEvent.Comments)).Scan(&id)

	//TODO update currentPosts when id > repo.currentEventId

	repo.currentEventId = id*/

	return
}

func (repo *repo) InsertFile(base64File string) (*model.File, error) {
	fileName := uuid.NewString()

	properties := strings.Split(base64File, ";base64,")
	contentType := strings.Split(properties[0], ":")

	uploadStream, err := repo.bucket.OpenUploadStream(
		fileName,
	)
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
	repo.currentPosts = append(repo.currentPosts, post)

	return post, nil
}

func (repo *repo) GetPosts() ([]*model.Post, error) {
	/*currentPosts := make([]*model.Post, 0)

	// first get all rows with event_type = "CreatePost" and latestEventId
	sqlQuery := `select "postId", description, data, liked, comments, id from "post-events" where id > $1 and "eventType" = $2 ORDER BY "id" ASC`

	rows, err := repo.client.Query(sqlQuery, repo.currentEventId, "CreatePost")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	oldId := repo.currentEventId
	id := repo.currentEventId

	for rows.Next() {
		var post model.Post
		var fileProperties string // "fileid.contenttype"

		err = rows.Scan(&post.ID, &post.Description, &fileProperties, pq.Array(&post.LikedBy), pq.Array(&post.Comments), &id)
		if err != nil {
			repo.currentEventId = oldId
			return nil, err
		}

		// id = 0, contenttype = 1
		idContentType := strings.Split(fileProperties, ".")

		content, err := util.LoadFile(idContentType[0])
		if err != nil {
			repo.currentEventId = oldId
			return nil, err
		}

		post.Data = &model.File{
			ID:          idContentType[0],
			Content:     content,
			ContentType: idContentType[1],
		}

		currentPosts = append(currentPosts, &post)
	}

	// when value of id hasn't changed, then list is already recent
	if id == repo.currentEventId {
		log.Printf("Post list is recent")
		return repo.currentPosts, nil
	}

	repo.currentEventId = id

	// then search for newest edit event (includes EditPost, LikePost, UnlikePost)
	// and add all data to posts
	for _, post := range currentPosts {
		sqlQuery = `select liked, description from "post-events" where id = (select max(id) from "post-events" where "postId" = $1 and ("eventType" = $2 or "eventType" = $3 or "eventType" = $4))`

		row := repo.client.QueryRow(sqlQuery, post.ID, "EditPost", "LikePost", "UnlikePost")

		switch err = row.Scan(pq.Array(&post.LikedBy), &post.Description); err {
		case sql.ErrNoRows:
			// nothing happens because it is not really an error
			// since a post doesn't have to be edited
		case nil:
			log.Printf("Edited data added to " + post.ID)
		default:
			repo.currentEventId = oldId
			return nil, err
		}
	}

	// Update runtime data
	repo.currentPosts = append(repo.currentPosts, currentPosts...)*/

	return repo.currentPosts, nil
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
