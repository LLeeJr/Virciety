package database

import (
	"database/sql"
	"posts-service/graph/model"
	"posts-service/util"
)

type PostEvent struct {
	EventTime   string    `json:"event_time"`
	EventType   string    `json:"event_type"`
	PostID      string    `json:"id"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	Data        string    `json:"data"`
	LikedBy     []*string `json:"liked_by"`
	Comments    []*string `json:"comments"`
}

type Repository interface {
	CreatePost(postEvent PostEvent) (*model.Post, error)
	GetPosts() ([]*model.Post, error)
	RemovePost(postEvent PostEvent) (string, error)
	EditPost(postEvent PostEvent) (string, error)
	LikePost(postEvent PostEvent) (string, error)
	UnlikePost(postEvent PostEvent) (string, error)
}

type repo struct {
	DB         *sql.DB
	PostEvents []*PostEvent
}

func NewRepo(db *sql.DB) (Repository, error) {
	return &repo{
		DB:         db,
		PostEvents: make([]*PostEvent, 0),
	}, nil
}

func (repo *repo) CreatePost(postEvent PostEvent) (*model.Post, error) {
	repo.PostEvents = append(repo.PostEvents, &postEvent)

	return &model.Post{
		ID:          postEvent.PostID,
		Description: postEvent.Description,
		Data:        postEvent.Data,
		LikedBy:     postEvent.LikedBy,
		Comments:    postEvent.Comments,
	}, nil
}

func (repo *repo) GetPosts( /*lastChecked string (eventTime)*/ ) ([]*model.Post, error) {
	currentPosts := make([]*model.Post, 0)

	// check from this last event timestamp, so we only get the recent events to process them

	// first get all rows with event_type = "PostCreated"

	for _, event := range repo.PostEvents {
		if event.EventType == "CreatePost" {
			post := &model.Post{
				ID:          event.PostID,
				Description: event.Description,
				Data:        event.Data,
				LikedBy:     event.LikedBy,
				Comments:    event.Comments,
			}

			currentPosts = append(currentPosts, post)
		}
	}

	for _, event := range repo.PostEvents {
		if event.EventType == "CreatePost" || event.EventType == "RemovePost" {
			continue
		}

		for _, post := range currentPosts {
			if event.PostID == post.ID {
				if event.EventType == "EditPost" {
					post.Description = event.Description
				} else if event.EventType == "LikePost" {
					post.LikedBy = append(post.LikedBy, event.LikedBy[len(event.LikedBy)-1])
				} else { // event.EventType == "UnlikePost"
					post.LikedBy = util.Compare(post.LikedBy, event.LikedBy)
				}

				break
			}
		}
	}

	// then get all rows with event_type = "PostUpdated"
	return currentPosts, nil
}

func (repo *repo) RemovePost(postEvent PostEvent) (string, error) {
	// delete all events relating to the id
	newPostEvents := make([]*PostEvent, 0)

	for _, event := range repo.PostEvents {
		if event.PostID == postEvent.PostID {
			continue
		}

		newPostEvents = append(newPostEvents, event)
	}

	newPostEvents = append(newPostEvents, &postEvent)

	// new current post events
	repo.PostEvents = newPostEvents

	return "success", nil
}

func (repo *repo) EditPost(postEvent PostEvent) (string, error) {
	repo.PostEvents = append(repo.PostEvents, &postEvent)

	return "success", nil
}

func (repo *repo) LikePost(postEvent PostEvent) (string, error) {
	repo.PostEvents = append(repo.PostEvents, &postEvent)

	return "success", nil
}

func (repo *repo) UnlikePost(postEvent PostEvent) (string, error) {
	repo.PostEvents = append(repo.PostEvents, &postEvent)

	return "success", nil
}
