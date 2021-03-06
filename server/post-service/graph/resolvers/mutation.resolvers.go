package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"post-service/database"
	"post-service/graph/generated"
	"post-service/graph/model"
	"time"
)

func (r *mutationResolver) CreatePost(ctx context.Context, newPost model.CreatePostRequest) (*model.Post, error) {
	// get current time
	created := time.Now().Format("2006-01-02 15:04:05")

	// create post event
	postEvent := database.PostEvent{
		EventTime:   created,
		EventType:   "CreatePost",
		Username:    newPost.Username,
		Description: newPost.Description,
		FileID:      "",
		LikedBy:     make([]string, 0),
		Comments:    make([]string, 0),
	}

	// save event in database
	post, err := r.repo.CreatePost(ctx, postEvent, newPost.Data)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	for _, observer := range r.observers {
		observer <- post
	}
	r.mu.Unlock()

	// put event on queue for notifications
	// r.producerQueue.AddMessageToQuery(postEvent)

	return post, nil
}

func (r *mutationResolver) EditPost(ctx context.Context, edit model.EditPostRequest) (string, error) {
	// process the data and create new post event
	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "EditPost",
		PostID:      edit.ID,
		Description: edit.NewDescription,
		LikedBy:     edit.LikedBy,
	}

	// save event in database
	ok, err := r.repo.EditPost(ctx, postEvent)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) RemovePost(ctx context.Context, remove model.RemovePostRequest) (string, error) {
	// process the data and create new post event
	postEvent := database.PostEvent{
		EventTime: time.Now().Format("2006-01-02 15:04:05"),
		EventType: "RemovePost",
		PostID:    remove.ID,
		FileID:    remove.FileID,
	}

	// save event in database
	ok, err := r.repo.RemovePost(ctx, postEvent)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) LikePost(ctx context.Context, like model.LikePostRequest) (string, error) {
	// process the data and create new post event
	event := "LikePost"

	if !like.Liked {
		event = "UnlikePost"
	}

	postEvent := database.PostEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   event,
		PostID:      like.ID,
		Username:    like.PostOwner,
		Description: like.Description,
		LikedBy:     like.NewLikedBy,
	}

	// save event in database
	_, err := r.repo.EditPost(ctx, postEvent)
	if err != nil {
		return "failed", err
	}

	// only send notification on like-post by different person than the post owner
	if like.Liked && like.PostOwner != like.LikedBy {
		// put event on queue for notifications
		r.producerQueue.AddMessageToEvent(postEvent)
	}

	return "success", nil
}

func (r *mutationResolver) AddComment(ctx context.Context, comment model.AddCommentRequest) (*model.Comment, error) {
	// create newComment
	newComment := &model.Comment{
		PostID:    comment.PostID,
		Comment:   comment.Comment,
		CreatedBy: comment.CreatedBy,
		Event:     "CreateComment",
	}

	// get post related to the comment
	post, err := r.repo.GetPost(ctx, comment.PostID)
	if err != nil {
		return nil, err
	}

	// create post comment event
	postCommentEvent := database.PostCommentEvent{
		Comment: newComment,
		Post:    post,
	}

	// put event on queue for comments
	r.producerQueue.AddMessageToCommand(postCommentEvent)

	return newComment, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
