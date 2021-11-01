package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"
	"user-service/database"
	"user-service/graph/generated"
	"user-service/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserData) (*model.User, error) {
	userEvent := database.UserEvent{
		EventType: "CreateUser",
		EventTime: time.Now(),
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Follows:   make([]string, 0),
	}

	user, err := r.repo.CreateUser(ctx, userEvent)

	if err != nil {
		return nil, err
	}

	r.publisher.AddMessageToEvent(userEvent, "User-Service")
	r.publisher.AddMessageToCommand("User-Service")
	//r.publisher.AddMessageToQuery()

	return user, nil
}

func (r *mutationResolver) AddFollow(ctx context.Context, id *string, toAdd *string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetUserByID(ctx context.Context, id *string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetUserByName(ctx context.Context, name *string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetFollows(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
