package database

import (
	"comment-service/graph/model"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type CommentEvent struct {
	EventTime   string   `json:"event_time"`
	EventType   string   `json:"event_type"`
	CommentID   string   `json:"id"`
	Username    string   `json:"username"`
	Description string   `json:"description"`
	LikedBy     []string `json:"liked_by"`
	PostID      string   `json:"post_id"`
}

type Repository interface {
	CreateComment(event CommentEvent) (*model.Comment, error)
	GetComments() (map[string][]*model.Comment, error)
	GetCurrentComments() map[string][]*model.Comment
	GetCommentsByPostId(postId string) ([]*model.Comment, error)
}

type repo struct {
	DB              *sql.DB
	currentComments map[string][]*model.Comment
	currentEventId  int
}

func NewRepo(db *sql.DB) (Repository, error) {
	return &repo{
		DB:              db,
		currentComments: make(map[string][]*model.Comment),
		currentEventId:  0,
	}, nil
}

func (repo *repo) InsertCommentEvent(commentEvent CommentEvent) (err error) {
	sqlQuery := `INSERT INTO "comment-events" ("commentID", "eventTime", "eventType", username, description, liked, "postID")
				VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`

	newTime, _ := time.Parse("2006-01-02 15:04:05", commentEvent.EventTime)

	id := 0

	err = repo.DB.QueryRow(sqlQuery, commentEvent.CommentID, newTime, commentEvent.EventType, commentEvent.Username, commentEvent.Description,
		pq.Array(commentEvent.LikedBy), commentEvent.PostID).Scan(&id)

	repo.currentEventId = id

	return
}

func (repo *repo) CreateComment(event CommentEvent) (*model.Comment, error) {
	err := repo.InsertCommentEvent(event)
	if err != nil {
		return nil, err
	}

	comment := &model.Comment{
		ID:          event.CommentID,
		Description: event.Description,
		LikedBy:     event.LikedBy,
		PostID:      event.PostID,
	}

	// add to currentComments
	comments, ok := repo.currentComments[comment.PostID]

	if ok {
		repo.currentComments[comment.PostID] = append(comments, comment)
	} else {
		repo.currentComments[comment.PostID] = []*model.Comment{comment}
	}

	return comment, nil
}

func (repo *repo) GetComments() (map[string][]*model.Comment, error) {
	currentComments := make([]*model.Comment, 0)

	// first get all rows with event_type = "CreateComment" and latestEventId
	sqlQuery := `select "commentID", description, liked, "postID", id from "comment-events" where id > $1 and "eventType" = $2 `

	rows, err := repo.DB.Query(sqlQuery, repo.currentEventId, "CreateComment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	oldId := repo.currentEventId
	id := repo.currentEventId

	for rows.Next() {
		var comment model.Comment

		err = rows.Scan(&comment.ID, &comment.Description, pq.Array(&comment.LikedBy), &comment.PostID, &id)
		if err != nil {
			repo.currentEventId = oldId
			return nil, err
		}
		currentComments = append(currentComments, &comment)
	}

	for _, comment := range currentComments {

		// add to currentComments
		comments, ok := repo.currentComments[comment.PostID]

		if ok {
			repo.currentComments[comment.PostID] = append(comments, comment)
		} else {
			repo.currentComments[comment.PostID] = []*model.Comment{comment}
		}
	}

	return repo.currentComments, nil
}

func (repo *repo) GetCurrentComments() map[string][]*model.Comment {
	return repo.currentComments
}

func (repo *repo) GetCommentsByPostId(postId string) ([]*model.Comment, error) {
	comments, ok := repo.currentComments[postId]
	if !ok {
		errMsg := "no comments for post with id " + postId + " found"
		return nil, errors.New(errMsg)
	}

	return comments, nil
}
