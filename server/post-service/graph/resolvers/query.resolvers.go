package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"posts-service/graph/generated"
	"posts-service/graph/model"
)

func (r *queryResolver) GetPosts(ctx context.Context, id string, fetchLimit int) ([]*model.Post, error) {
	currentPosts := r.repo.GetCurrentPosts()

	fmt.Printf("len(currentPosts): %d\n", len(currentPosts))
	fmt.Printf("id: %s\n", id)
	fmt.Printf("fetchLimit: %d\n", fetchLimit)

	// cached backend data isn't enough, so fetch rest from db
	if len(currentPosts) < fetchLimit && id == "" {
		newPosts, err := r.repo.GetPosts(fetchLimit - len(currentPosts))
		if err != nil {
			return nil, err
		}
		currentPosts = append(currentPosts, newPosts...)
		// cached backend data is enough
	} else if len(currentPosts) == fetchLimit && id == "" {
		return currentPosts, nil
		// cached backend data is more than enough
	} else if len(currentPosts) > fetchLimit && id == "" {
		return currentPosts[:fetchLimit], nil
		// check how much data need to be fetched from db
	} else if len(currentPosts) <= fetchLimit && id != "" {
		// get index of post with id and adjust the fetchLimit
		index, _ := r.repo.GetPostById(id)
		newFetchLimit := fetchLimit - (len(currentPosts) - (index + 1))
		if index == -1 {
			return nil, errors.New("post with id " + id + " not found while fetching")
		}
		newPosts, err := r.repo.GetPosts(newFetchLimit)
		if err != nil {
			return nil, err
		}
		currentPosts = append(currentPosts[index+1:], newPosts...)
		// check if data needs to be fetched from db or remaining data from cached backend data is enough
	} else if len(currentPosts) > fetchLimit && id != "" {
		index, _ := r.repo.GetPostById(id)
		if index == -1 {
			return nil, errors.New("post with id " + id + " not found while fetching")
		}
		// cached backend data is enough
		if len(currentPosts)-(index+1) > fetchLimit {
			return currentPosts[index+1 : index+1+fetchLimit], nil
			// cached backend data is just enough
		} else if len(currentPosts)-(index+1) == fetchLimit {
			return currentPosts[index+1:], nil
			// data needs to be fetched from db
		} else {
			newPosts, err := r.repo.GetPosts(fetchLimit - (len(currentPosts) - (index + 1)))
			if err != nil {
				return nil, err
			}
			currentPosts = append(currentPosts[index+1:], newPosts...)
		}
	}

	return currentPosts, nil
}

func (r *queryResolver) GetData(ctx context.Context, id string) (string, error) {
	currentPosts := r.repo.GetCurrentPosts()

	for _, post := range currentPosts {
		if post.ID == id {
			return post.Data.Content, nil
		}
	}

	return "", errors.New("post with id " + id + " not found")
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
