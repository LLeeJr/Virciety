package database

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"posts-service/graph/model"
	"time"
)

type Repository interface {
	InsertPostEvent(post PostEvent) error
	CreatePost(postEvent PostEvent) (*model.Post, error)
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
	DB             *sql.DB
	currentPosts   []*model.Post
	currentEventId int
}

func NewRepo(db *sql.DB) (Repository, error) {
	return &repo{
		DB:             db,
		currentEventId: 0,
		currentPosts:   make([]*model.Post, 0),
	}, nil
}

func (repo *repo) InsertPostEvent(postEvent PostEvent) (err error) {
	sqlQuery := `INSERT INTO "post-events" ("postId", "eventTime", "eventType", username, description, data, liked, comments)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	newTime, _ := time.Parse("2006-01-02 15:04:05", postEvent.EventTime)

	id := 0

	err = repo.DB.QueryRow(sqlQuery, postEvent.PostID, newTime, postEvent.EventType, postEvent.Username, postEvent.Description, postEvent.Data,
		pq.Array(postEvent.LikedBy), pq.Array(postEvent.Comments)).Scan(&id)

	repo.currentEventId = id

	return
}

func (repo *repo) CreatePost(postEvent PostEvent) (*model.Post, error) {
	err := repo.InsertPostEvent(postEvent)
	if err != nil {
		return nil, err
	}

	post := &model.Post{
		ID:          postEvent.PostID,
		Description: postEvent.Description,
		Data:        postEvent.Data,
		LikedBy:     postEvent.LikedBy,
		Comments:    postEvent.Comments,
	}

	// add to currentPosts
	repo.currentPosts = append(repo.currentPosts, post)

	return post, nil
}

func (repo *repo) GetPosts() ([]*model.Post, error) {
	currentPosts := make([]*model.Post, 0)

	// first get all rows with event_type = "CreatePost" and latestEventId
	sqlQuery := `select "postId", description, data, liked, comments, id from "post-events" where id > $1 and "eventType" = $2 `

	rows, err := repo.DB.Query(sqlQuery, repo.currentEventId, "CreatePost")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	oldId := repo.currentEventId
	id := repo.currentEventId

	for rows.Next() {
		var post model.Post

		err = rows.Scan(&post.ID, &post.Description, &post.Data, pq.Array(&post.LikedBy), pq.Array(&post.Comments), &id)
		if err != nil {
			repo.currentEventId = oldId
			return nil, err
		}
		currentPosts = append(currentPosts, &post)
	}

	// list is recent
	if id == repo.currentEventId {
		return repo.currentPosts, nil
	}

	repo.currentEventId = id

	// then search for newest edit event (includes EditPost, LikePost, UnlikePost)
	// and add all data to posts
	for _, post := range currentPosts {
		sqlQuery = `select liked, description from "post-events" where id = (select max(id) from "post-events" where "postId" = $1 and ("eventType" = $2 or "eventType" = $3 or "eventType" = $4))`

		row := repo.DB.QueryRow(sqlQuery, post.ID, "EditPost", "LikePost", "UnlikePost")

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
	repo.currentPosts = append(repo.currentPosts, currentPosts...)

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
	// remove from currentPosts
	repo.currentPosts = append(repo.currentPosts[:index], repo.currentPosts[index+1:]...)

	// delete all events relating to the id
	sqlQuery := `delete from "post-events" where "postId" = $1`

	_, err := repo.DB.Exec(sqlQuery, postEvent.PostID)
	if err != nil {
		return "failed", err
	}

	// new current post events
	err = repo.InsertPostEvent(postEvent)
	if err != nil {
		return "failed", err
	}

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
