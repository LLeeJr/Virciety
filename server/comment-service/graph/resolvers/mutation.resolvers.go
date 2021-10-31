package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"comment-service/graph/generated"
	"comment-service/graph/model"
	"context"
)

func (r *mutationResolver) CreateComment(ctx context.Context, newComment model.CreateCommentRequest) (*model.Comment, error) {
	/*created := time.Now().Format("2006-01-02 15:04:05")

	// create new event
	commentEvent := database.CommentEvent{
		EventTime:   created,
		EventType:   "CreateComment",
		CommentID:   created + "__" + newComment.Username + "__" + newComment.PostID,
		PostID:      newComment.PostID,
	}

	// save event in database
	comment, err := r.repo.CreateComment(commentEvent)
	if err != nil {
		return nil, err
	}

	// put event on queue for notifications
	// put event on queue for posts
	r.producerQueue.AddMessageToEvent(commentEvent)*/

	return nil, nil
}

func (r *mutationResolver) EditComment(ctx context.Context, edit model.EditCommentRequest) (string, error) {
	/*// get comment by id
	comment, _, username, err := r.repo.GetCommentById(edit.ID)
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
		Username:    username,
		Description: comment.Description,
		LikedBy:     comment.LikedBy,
	}

	// save event in database
	ok, err := r.repo.EditComment(commentEvent)
	if err != nil {
		return ok, err
	}*/

	return "", nil
}

func (r *mutationResolver) RemoveComment(ctx context.Context, removeID string) (string, error) {
	/*// get comment by id
	comment, index, username, err := r.repo.GetCommentById(removeID)
	if err != nil {
		return "failed", err
	}

	// add new data and create event
	commentEvent := database.CommentEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "RemoveComment",
		CommentID:   comment.ID,
		PostID:      comment.PostID,
		Username:    username,
		Description: comment.Description,
		LikedBy:     make([]string, 0),
	}

	// save event in database
	ok, err := r.repo.RemoveComment(commentEvent, index)
	if err != nil {
		return ok, err
	}*/

	return "", nil
}

func (r *mutationResolver) LikeComment(ctx context.Context, like model.UnLikeCommentRequest) (string, error) {
	/*// get comment by id
	comment, _, username, err := r.repo.GetCommentById(like.ID)
	if err != nil {
		return "failed", err
	}

	// add new data and create event
	if util.Contains(comment.LikedBy, like.Username) {
		errMsg := "user " + like.Username + " already liked the comment with id " + like.ID
		return "failed", errors.New(errMsg)
	}

	comment.LikedBy = append(comment.LikedBy, like.Username)

	commentEvent := database.CommentEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "LikeComment",
		CommentID:   comment.ID,
		PostID:      comment.PostID,
		Username:    username,
		Description: comment.Description,
		LikedBy:     comment.LikedBy,
	}

	// save event in database
	ok, err := r.repo.LikeComment(commentEvent)
	if err != nil {
		return ok, err
	}*/

	return "", nil
}

func (r *mutationResolver) UnlikeComment(ctx context.Context, unlike model.UnLikeCommentRequest) (string, error) {
	/*// get comment by id
	comment, _, username, err := r.repo.GetCommentById(unlike.ID)
	if err != nil {
		return "failed", err
	}

	// add new data and create event
	index := util.Search(comment.LikedBy, unlike.Username)

	if index == -1 {
		errMsg := "user " + unlike.Username + " hasn't liked the comment with id " + unlike.ID
		return "failed", errors.New(errMsg)
	}

	comment.LikedBy = append(comment.LikedBy[:index], comment.LikedBy[index+1:]...)

	commentEvent := database.CommentEvent{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventType:   "UnlikeComment",
		CommentID:   comment.ID,
		PostID:      comment.PostID,
		Username:    username,
		Description: comment.Description,
		LikedBy:     comment.LikedBy,
	}

	// save event in database
	ok, err := r.repo.UnlikeComment(commentEvent)
	if err != nil {
		return ok, err
	}*/

	return "", nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
