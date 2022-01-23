package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"
	"post-service/graph/generated"
	"post-service/graph/model"
	"post-service/util"

	"github.com/google/uuid"
)

func (r *queryResolver) GetPosts(ctx context.Context, id string, fetchLimit int, filter *string) ([]*model.Post, error) {
	if filter == nil {
		log.Println("Filter == nil")
	} else {
		log.Printf("Filter != nil: %s", *filter)
	}

	currentPosts, err := r.repo.GetPosts(id, fetchLimit, filter)
	if err != nil {
		return nil, err
	}

	return currentPosts, nil
}

func (r *queryResolver) GetData(ctx context.Context, fileID string) (string, error) {
	data, err := r.repo.GetData(fileID)
	if err != nil {
		return "", err
	}

	return data, nil
}

func (r *queryResolver) GetPostComments(ctx context.Context, id string) (*model.CommentsWithProfileIds, error) {
	requestID := uuid.NewString()

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(r.responses, requestID)
		r.mu.Unlock()
	}()

	r.mu.Lock()
	r.responses[requestID] = make(chan []*model.Comment, 1)
	r.mu.Unlock()

	r.producerQueue.AddMessageToQuery(id, requestID)

	comments := <-r.responses[requestID]

	users := make([]string, 0)
	for _, comment := range comments {
		if !util.Contains(users, comment.CreatedBy) {
			users = append(users, comment.CreatedBy)
		}
	}

	r.mu.Lock()
	r.userResponses[requestID] = make(chan map[string]string, 1)
	r.mu.Unlock()

	r.producerQueue.ProfilePictureIdQuery(id, requestID, users)

	userIdMap := <-r.userResponses[requestID]
	m := make([]*model.UserIDMap, 0)
	for user, pictureId := range userIdMap {
		log.Printf(user, pictureId)
		m = append(m, &model.UserIDMap{
			Key:   user,
			Value: pictureId,
		})
	}

	commentIdMap := &model.CommentsWithProfileIds{
		Comments:  comments,
		UserIDMap: m,
	}

	return commentIdMap, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
