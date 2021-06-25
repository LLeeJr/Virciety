package transport

import "posts-service/service"

type (
	CreatePostRequest struct {
		Username 	string 	`json:"username"`
		Description string 	`json:"description"`
		Data 		string 	`json:"data"`
	}
	CreatePostResponse struct {
		Post service.Post `json:"post"`
		Err  string       `json:"error,omitempty"`
	}

	GetPostsRequest struct {}
	GetPostsResponse struct {
		Posts 		[]service.Post `json:"posts"`
		Err			string       `json:"error,omitempty"`
	}

	RemovePostRequest struct {
		ID 			string 	`json:"id"`
	}
	RemovePostResponse struct {
		Removed		bool	`json:"removed"`
		Err			string	`json:"error,omitempty"`
	}

	EditPostRequest struct {
		ID			string  `json:"id"`
		Description string  `json:"description"`
	}
	EditPostResponse struct {
		Updated		bool	`json:"updated"`
		Err			string	`json:"error,omitempty"`
	}

	LikedPostRequest struct {
		ID			string	`json:"id"`
		Username	string	`json:"username"`
	}
	LikedPostResponse struct {
		Liked		bool	`json:"liked"`
		Err			string	`json:"error,omitempty"`
	}
)
