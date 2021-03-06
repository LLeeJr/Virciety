// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Comment struct {
	ID          string   `json:"id"`
	PostID      string   `json:"post_id"`
	Description string   `json:"description"`
	LikedBy     []string `json:"likedBy"`
}

type CreateCommentRequest struct {
	Username    string `json:"username"`
	Description string `json:"description"`
	PostID      string `json:"postID"`
}

type EditCommentRequest struct {
	ID             string `json:"id"`
	NewDescription string `json:"newDescription"`
}

type MapComments struct {
	Key   string     `json:"key"`
	Value []*Comment `json:"value"`
}

type RemoveCommentRequest struct {
	ID string `json:"id"`
}

type UnLikeCommentRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
