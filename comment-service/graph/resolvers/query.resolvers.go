package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"comment-service/graph/generated"
	"comment-service/graph/model"
	"comment-service/util"
	"context"
)

func (r *queryResolver) GetComments(ctx context.Context) ([]*model.MapComments, error) {
	currentComments, err := r.repo.GetComments()
	if err != nil {
		return nil, err
	}

	return util.ConvertedIntoMapComments(currentComments), nil
}

func (r *queryResolver) GetCommentsByPostID(ctx context.Context, id string) ([]*model.Comment, error) {
	_, err := r.repo.GetComments()
	if err != nil {
		return nil, err
	}

	comments, err := r.repo.GetCommentsByPostId(id)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
