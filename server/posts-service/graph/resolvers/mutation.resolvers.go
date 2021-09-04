package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"posts-service/database"
	"posts-service/graph/generated"
	"posts-service/graph/model"
	"posts-service/util"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (r *mutationResolver) CreatePost(ctx context.Context, newPost model.CreatePostRequest) (*model.Post, error) {
	created := time.Now().Format("2006-01-02 15:04:05")

	properties := strings.Split(newPost.Data, ";base64,")
	contentType := strings.Split(properties[0], ":")

	postEvent := database.PostEvent{
		EventTime:   created,
		EventType:   "CreatePost",
		PostID:      created + "__" + newPost.Username,
		Username:    newPost.Username,
		Description: newPost.Description,
		Data: &model.File{
			ID:          uuid.NewString(),
			Content:     properties[1],
			ContentType: contentType[1],
		},
		LikedBy:  make([]string, 0),
		Comments: make([]string, 0),
	}

	// save event in database
	post, err := r.repo.CreatePost(postEvent)
	if err != nil {
		return nil, err
	}

	// put event on queue for notifications
	// r.producerQueue.AddMessageToQuery(postEvent)

	return post, nil
}

func (r *mutationResolver) Upload(ctx context.Context, file string) (*model.File, error) {
	properties := strings.Split(file, ";base64,")
	contentType := strings.Split(properties[0], ":")

	return &model.File{
		ID:          uuid.NewString(),
		Content:     properties[1],
		ContentType: contentType[1],
	}, nil
}

func (r *mutationResolver) EditPost(ctx context.Context, edit model.EditPostRequest) (string, error) {
	// get post data out of current saved posts and remove it out of the list
	_, post := r.repo.GetPostById(edit.ID)
	if post == nil {
		errMsg := "no post with id " + edit.ID + " found"
		return "failed", errors.New(errMsg)
	}

	// process the data and create new post event
	info := strings.Split(post.ID, "__")

	post.Description = edit.NewDescription

	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "EditPost",
		PostID:      post.ID,
		Username:    info[1],
		Description: post.Description,
		Data:        post.Data,
		LikedBy:     post.LikedBy,
		Comments:    post.Comments,
	}

	// save event in database
	ok, err := r.repo.EditPost(postEvent)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) RemovePost(ctx context.Context, removeID string) (string, error) {
	// get post data out of current saved posts
	index, post := r.repo.GetPostById(removeID)
	if post == nil {
		errMsg := "no post with id " + removeID + " found"
		return "failed", errors.New(errMsg)
	}

	// process the data and create new post event
	info := strings.Split(post.ID, "__")

	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "RemovePost",
		PostID:      post.ID,
		Username:    info[1],
		Description: post.Description,
		Data:        post.Data,
		LikedBy:     make([]string, 0),
		Comments:    make([]string, 0),
	}

	// save event in database
	ok, err := r.repo.RemovePost(postEvent, index)
	if err != nil {
		return ok, err
	}

	// put event on queue for notifications to remove all notification events for this post
	// put event on queue for comments to remove all comment events for this post
	// r.producerQueue.AddMessageToQuery(postEvent)

	return ok, nil
}

func (r *mutationResolver) LikePost(ctx context.Context, like model.UnLikePostRequest) (string, error) {
	// get post data out of current saved posts and remove it out of the list
	_, post := r.repo.GetPostById(like.ID)
	if post == nil {
		errMsg := "no post with id " + like.ID + " found"
		return "failed", errors.New(errMsg)
	}

	// process the data and create new post event
	info := strings.Split(post.ID, "__")

	if util.Contains(post.LikedBy, like.Username) {
		errMsg := "user " + like.Username + " already liked the post with id " + like.ID
		return "failed", errors.New(errMsg)
	}

	post.LikedBy = append(post.LikedBy, like.Username)

	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "LikePost",
		PostID:      post.ID,
		Username:    info[1],
		Description: post.Description,
		Data:        post.Data,
		LikedBy:     post.LikedBy,
		Comments:    post.Comments,
	}

	// save event in database
	ok, err := r.repo.LikePost(postEvent)
	if err != nil {
		return ok, err
	}

	// put event on queue for notifications
	// r.producerQueue.AddMessageToQuery(postEvent)

	return ok, nil
}

func (r *mutationResolver) UnlikePost(ctx context.Context, unlike model.UnLikePostRequest) (string, error) {
	// get post data out of current saved posts and remove it out of the list
	_, post := r.repo.GetPostById(unlike.ID)
	if post == nil {
		errMsg := "no post with id " + unlike.ID + " found"
		return "failed", errors.New(errMsg)
	}

	// process the data and create new post event
	info := strings.Split(post.ID, "__")

	index := util.Search(post.LikedBy, unlike.Username)

	if index == -1 {
		errMsg := "user " + unlike.Username + " hasn't liked the post with id " + unlike.ID
		return "failed", errors.New(errMsg)
	}

	post.LikedBy = append(post.LikedBy[:index], post.LikedBy[index+1:]...)

	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "UnlikePost",
		PostID:      post.ID,
		Username:    info[1],
		Description: post.Description,
		Data:        post.Data,
		LikedBy:     post.LikedBy,
		Comments:    post.Comments,
	}

	// save event in database
	ok, err := r.repo.UnlikePost(postEvent)
	if err != nil {
		return ok, err
	}

	// put event on queue for notifications for deleting like notification?
	// r.producerQueue.AddMessageToQuery(postEvent)

	return ok, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
