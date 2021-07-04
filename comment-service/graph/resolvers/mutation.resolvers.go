package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"comment-service/database"
	"comment-service/graph/generated"
	"comment-service/graph/model"
	"context"
	"fmt"
	"time"
)

func (r *mutationResolver) CreateComment(ctx context.Context, newComment model.CreateCommentRequest) (*model.Comment, error) {
	created := time.Now().Format("2006-01-02 15:04:05")

	// create new event
	commentEvent := database.CommentEvent{
		EventTime:   created,
		EventType:   "CreateComment",
		CommentID:   created + "__" + newComment.Username + "__" + newComment.PostID,
		PostID:      newComment.PostID,
		Username:    newComment.Username,
		Description: newComment.Description,
		LikedBy:     make([]string, 0),
	}

	// save event in database
	comment, err := r.repo.CreateComment(commentEvent)
	if err != nil {
		return nil, err
	}

	// put event on queue for notifications
	// put event on queue for posts
	r.producerQueue.AddMessageToEvent(commentEvent, "Post-Service")

	return comment, nil
}

func (r *mutationResolver) EditComment(ctx context.Context, edit model.EditCommentRequest) (string, error) {
	// get comment by id
	comment, postId, err := r.repo.GetCommentById(edit.ID)
	if err != nil {
		return "failed", err
	}

	// add new data and create event
	comment.Description = edit.NewDescription

	commentEvent := database.CommentEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "EditComment",
		CommentID:   comment.ID,
		PostID:      comment.PostID,
		Username:    postId,
		Description: comment.Description,
		LikedBy:     comment.LikedBy,
	}

	// save event in database
	ok, err := r.repo.EditComment(commentEvent)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (r *mutationResolver) RemoveComment(ctx context.Context, removeID string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) LikeComment(ctx context.Context, like model.UnLikeCommentRequest) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnlikeComment(ctx context.Context, unlike model.UnLikeCommentRequest) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
