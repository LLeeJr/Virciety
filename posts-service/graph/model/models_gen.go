// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreatePostRequest struct {
	Username    string `json:"username"`
	Description string `json:"description"`
	Data        string `json:"data"`
}

type EditPostRequest struct {
	ID             string `json:"id"`
	NewDescription string `json:"newDescription"`
}

type Post struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Data        string    `json:"data"`
	LikedBy     []*string `json:"likedBy"`
	Comments    []*string `json:"comments"`
}

type RemovePostRequest struct {
	ID string `json:"id"`
}

type UnLikePostRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
