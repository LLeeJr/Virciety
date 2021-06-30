package database

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"posts-service/graph/model"
	"time"
)

type PostEvent struct {
	EventTime   string   `json:"event_time"`
	EventType   string   `json:"event_type"`
	PostID      string   `json:"id"`
	Username    string   `json:"username"`
	Description string   `json:"description"`
	Data        string   `json:"data"`
	LikedBy     []string `json:"liked_by"`
	Comments    []string `json:"comments"`
}

type Repository interface {
	InsertPostEvent(post PostEvent) error
	CreatePost(postEvent PostEvent) (*model.Post, error)
	GetPosts() ([]*model.Post, error)
	RemovePost(postEvent PostEvent) (string, error)
	EditPost(postEvent PostEvent) (string, error)
	LikePost(postEvent PostEvent) (string, error)
	UnlikePost(postEvent PostEvent) (string, error)
}

type repo struct {
	DB             *sql.DB
	currentEventId int
}

func NewRepo(db *sql.DB) (Repository, error) {
	return &repo{
		DB:             db,
		currentEventId: 0,
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

	return &model.Post{
		ID:          postEvent.PostID,
		Description: postEvent.Description,
		Data:        postEvent.Data,
		LikedBy:     postEvent.LikedBy,
		Comments:    postEvent.Comments,
	}, nil
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
		return nil, nil
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

	return currentPosts, nil
}

func (repo *repo) RemovePost(postEvent PostEvent) (string, error) {
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
