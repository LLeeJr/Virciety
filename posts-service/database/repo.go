package database

import (
	"database/sql"
	"posts-service/graph/model"
)

type PostEvent struct {
	EventTime 	string 		`json:"event_time"`
	EventType 	string 		`json:"event_type"`
	Username 	string 		`json:"username"`
	Description string 		`json:"description"`
	Data 		string 		`json:"data"`
	LikedBy		[]string 	`json:"liked_by"`
	Comments	[]string 	`json:"comments"`
}

type Repository interface {
	CreatePost(postEvent PostEvent) (*model.Post, error)
	GetPosts() ([]*model.Post, error)
}

type repo struct {
	DB         *sql.DB
	PostEvents []PostEvent
}

func NewRepo(db *sql.DB) (Repository, error) {
	return &repo{PostEvents: make([]PostEvent, 0)}, nil
}

func (repo *repo) CreatePost(postEvent PostEvent) (*model.Post, error) {
	repo.PostEvents = append(repo.PostEvents, postEvent)

	 return &model.Post{
		ID:          postEvent.EventTime + postEvent.Username,
		Description: postEvent.Description,
		Data:        postEvent.Data,
		LikedBy:     postEvent.LikedBy,
		Comments:    postEvent.Comments,
	}, nil
}

func (repo *repo) GetPosts() ([]*model.Post, error) {
	currentPosts := make([]*model.Post, 0)

	// first get all rows with event_type = "PostCreated"
	for _, event := range repo.PostEvents {
		post := &model.Post{
			ID:          event.EventTime + event.Username,
			Description: event.Description,
			Data:        event.Data,
			LikedBy:     event.LikedBy,
			Comments:    event.Comments,
		}

		currentPosts = append(currentPosts, post)
	}

	return currentPosts, nil

	// then get all rows with event_type = "PostUpdated"
}