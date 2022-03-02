package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
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
		Followers: make([]string, 0),
		FileId:    "",
	}

	user, err := r.repo.CreateUser(ctx, userEvent)

	if err != nil {
		return nil, err
	}

	//r.publisher.AddMessageToEvent(userEvent, "User-Service")
	r.publisher.AddMessageToCommand("User-Service")
	//r.publisher.AddMessageToQuery()

	return user, nil
}

func (r *mutationResolver) AddFollow(ctx context.Context, id *string, username *string, toAdd *string) (*model.User, error) {
	user, err := r.repo.AddFollow(ctx, id, username, toAdd)

	if err != nil {
		return nil, err
	}

	followEvent := &database.FollowEvent{
		EventType: "New Follower",
		EventTime: time.Now(),
		Username:  *toAdd,
		NewFollower: *username,
	}

	r.publisher.AddMessageToEvent(followEvent, "User-Service")

	return user, nil
}

func (r *mutationResolver) RemoveFollow(ctx context.Context, id *string, username *string, toRemove *string) (*model.User, error) {
	user, err := r.repo.RemoveFollow(ctx, id, username, toRemove)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) AddProfilePicture(ctx context.Context, input model.AddProfilePicture) (string, error) {
	profilePictureEvent := database.ProfilePictureEvent{
		EventType: "AddProfilePicture",
		EventTime: time.Now(),
		Username:  input.Username,
		FileId:    "",
	}

	response, err := r.repo.AddProfilePicture(ctx, profilePictureEvent, input)
	if err != nil {
		return "error while adding profile picture", err
	}

	return response, nil
}

func (r *mutationResolver) RemoveProfilePicture(ctx context.Context, remove model.RemoveProfilePicture) (string, error) {
	profilePictureEvent := database.ProfilePictureEvent{
		EventType: "RemoveProfilePicture",
		EventTime: time.Now(),
		Username:  remove.Username,
		FileId:    remove.FileID,
	}

	response, err := r.repo.RemoveProfilePicture(ctx, profilePictureEvent)
	if err != nil {
		return "error while removing profile picture", err
	}

	return response, nil
}

func (r *queryResolver) GetUserByID(ctx context.Context, id *string) (*model.User, error) {
	user, err := r.repo.GetUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *queryResolver) GetUserByName(ctx context.Context, name *string) (*model.User, error) {
	user, err := r.repo.GetUserByName(ctx, name)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *queryResolver) GetProfilePicture(ctx context.Context, fileID *string) (string, error) {
	data, err := r.repo.GetProfilePicture(ctx, fileID)
	if err != nil {
		return "error while retrieving profile picture", err
	}

	return data, nil
}

func (r *queryResolver) FindUsersWithName(ctx context.Context, name *string) ([]*model.User, error) {
	users, err := r.repo.FindUsersWithName(ctx, name)

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		msg := fmt.Sprint("no users found with given name: ", *name)
		return nil, errors.New(msg)
	}

	return users, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
