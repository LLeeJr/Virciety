package transport

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
	"posts-service/service"
)

func MakeCreatePostEndpoint(svc service.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreatePostRequest)
		post, err := svc.CreatePost(ctx, req.Username, req.Description, req.Data)
		if err != nil {
			return CreatePostResponse{
				Post: post,
				Err:  err.Error(),
			}, nil
		}
		return CreatePostResponse{
			Post: post,
			Err:  "",
		}, nil
	}
}

func MakeGetPostsEndpoint(svc service.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		posts, err := svc.GetPosts(ctx)
		if err != nil {
			return GetPostsResponse{
				Posts: nil,
				Err:   "no data found",
			}, nil
		}
		return GetPostsResponse{
			Posts: posts,
			Err:   "",
		}, nil
	}
}

func MakeRemovePostEndpoint(svc service.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RemovePostRequest)
		removed, err := svc.RemovePost(ctx, req.ID)
		if err != nil {
			return RemovePostResponse{
				Removed: removed,
				Err:     err.Error(),
			}, nil
		}
		return RemovePostResponse{
			Removed: removed,
			Err:     "",
		}, nil
	}
}

func MakeEditPostEndpoint(svc service.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EditPostRequest)
		updated, err := svc.EditPost(ctx, req.ID, req.Description)
		if err != nil {
			return EditPostResponse{
				Updated: updated,
				Err:  err.Error(),
			}, nil
		}
		return EditPostResponse{
			Updated: updated,
			Err:  "",
		}, nil
	}
}

func MakeLikedPostEndpoint(svc service.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LikedPostRequest)
		liked, err := svc.LikedPost(ctx, req.ID, req.Username)
		if err != nil {
			return LikedPostResponse{
				Liked: liked,
				Err:   err.Error(),
			}, nil
		}
		return LikedPostResponse{
			Liked: liked,
			Err:   "",
		}, nil
	}
}

func DecodeCreatePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeGetPostsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req GetPostsRequest
	return req, nil
}

func DecodeRemovePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req RemovePostRequest
	vars := mux.Vars(r)
	req = RemovePostRequest{
		ID: vars["id"],
	}
	return req, nil
}

func DecodeEditPostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req EditPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeLikedPostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req LikedPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}