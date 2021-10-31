package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/google/uuid"
	"posts-service/graph/generated"
	"posts-service/graph/model"
)

func (r *queryResolver) GetPosts(ctx context.Context, id string, fetchLimit int) ([]*model.Post, error) {
	currentPosts, err := r.repo.GetPosts(id, fetchLimit)
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

func (r *queryResolver) GetPostComments(ctx context.Context, id string) ([]*model.Comment, error) {
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

	return comments, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
