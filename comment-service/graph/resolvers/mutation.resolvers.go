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
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
