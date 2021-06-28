package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"posts-service/database"
	"posts-service/graph/generated"
	"posts-service/graph/model"
	"strings"
	"time"
)

func (r *mutationResolver) CreatePost(_ context.Context, newPost model.CreatePostRequest) (*model.Post, error) {
	created := time.Now().String()

	postEvent := database.PostEvent{
		EventTime:   created,
		EventType:   "CreatePost",
		PostID:      created + "__" + newPost.Username,
		Username:    newPost.Username,
		Description: newPost.Description,
		Data:        newPost.Data,
		LikedBy:     make([]*string, 0),
		Comments:    make([]*string, 0),
	}

	// put event on queue for notifs?

	// save event in database
	post, err := r.repo.CreatePost(postEvent)
	if err != nil {
		return nil, err
	}

	// add to currentPosts in resolver
	r.currentPosts = append(r.currentPosts, post)

	return post, nil
}

func (r *mutationResolver) EditPost(ctx context.Context, edit model.EditPostRequest) (string, error) {
	// get post data out of current saved posts and remove it out of the list
	_, post := GetPostByID(r.currentPosts, edit.ID)
	if post == nil {
		errMsg := "no post with id " + edit.ID + " found"
		return "failed", errors.New(errMsg)
	}

	// process the data and create new post event
	info := strings.Split(post.ID, "__")

	post.Description = edit.NewDescription

	postEvent := database.PostEvent{
		EventTime:   time.Now().String(),
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
	// get post data out of current saved posts and remove it out of the list
	index, post := GetPostByID(r.currentPosts, removeID)
	if post == nil {
		errMsg := "no post with id " + removeID + " found"
		return "failed", errors.New(errMsg)
	}

	r.currentPosts = append(r.currentPosts[:index], r.currentPosts[index+1:]...)

	// process the data and create new post event
	info := strings.Split(post.ID, "__")

	postEvent := database.PostEvent{
		EventTime:   time.Now().String(),
		EventType:   "RemovePost",
		PostID:      post.ID,
		Username:    info[1],
		Description: post.Description,
		Data:        post.Data,
		LikedBy:     make([]*string, 0),
		Comments:    make([]*string, 0),
	}

	// save event in database
	ok, err := r.repo.RemovePost(postEvent)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) LikePost(ctx context.Context, like model.UnLikePostRequest) (string, error) {
	// get post data out of current saved posts and remove it out of the list
	_, post := GetPostByID(r.currentPosts, like.ID)
	if post == nil {
		errMsg := "no post with id " + like.ID + " found"
		return "failed", errors.New(errMsg)
	}

	// process the data and create new post event
	info := strings.Split(post.ID, "__")

	if contains(post.LikedBy, like.Username) {
		errMsg := "user " + like.Username + " already liked the post with id " + like.ID
		return "failed", errors.New(errMsg)
	}

	post.LikedBy = append(post.LikedBy, &like.Username)

	postEvent := database.PostEvent{
		EventTime:   time.Now().String(),
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

	return ok, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
