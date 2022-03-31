package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"post-service/graph/generated"
	"post-service/graph/model"
	"post-service/util"

	"github.com/google/uuid"
)

func (r *queryResolver) GetPosts(ctx context.Context, id string, fetchLimit int, filter *string) ([]*model.Post, error) {
	// get posts from database
	currentPosts, err := r.repo.GetPosts(ctx, id, fetchLimit, filter)
	if err != nil {
		return nil, err
	}

	return currentPosts, nil
}

func (r *queryResolver) GetData(ctx context.Context, fileID string) (string, error) {
	// get file data from database
	data, err := r.repo.GetData(ctx, fileID)
	if err != nil {
		return "", err
	}

	return data, nil
}

func (r *queryResolver) GetPostComments(ctx context.Context, id string) (*model.CommentsWithProfileIds, error) {
	// set a new request id
	requestID := uuid.NewString()

	// delete request when context is done
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(r.responses, requestID)
		r.mu.Unlock()
	}()

	// add request id to comment response map
	r.mu.Lock()
	r.responses[requestID] = make(chan []*model.Comment, 1)
	r.mu.Unlock()

	// put request on queue for comment service
	r.producerQueue.AddMessageToQuery(id, requestID)

	// wait for response
	comments := <-r.responses[requestID]

	// get unique users from comments
	users := make([]string, 0)
	for _, comment := range comments {
		if !util.Contains(users, comment.CreatedBy) {
			users = append(users, comment.CreatedBy)
		}
	}

	// add request id to user response map
	r.mu.Lock()
	r.userResponses[requestID] = make(chan map[string]string, 1)
	r.mu.Unlock()

	// put request on queue for user service
	r.producerQueue.ProfilePictureIdQuery(id, requestID, users)

	// wait for response
	userIdMap := <-r.userResponses[requestID]

	// map user and picture id
	m := make([]*model.UserIDMap, 0)
	for user, pictureId := range userIdMap {
		m = append(m, &model.UserIDMap{
			Key:   user,
			Value: pictureId,
		})
	}

	// map comment and user id
	commentIdMap := &model.CommentsWithProfileIds{
		Comments:  comments,
		UserIDMap: m,
	}

	return commentIdMap, nil
}

func (r *queryResolver) GetPost(ctx context.Context, id string) (*model.Post, error) {
	// get post by given id
	post, err := r.repo.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}

	// get file data
	post.Data.Content, err = r.GetData(ctx, post.Data.Name)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
