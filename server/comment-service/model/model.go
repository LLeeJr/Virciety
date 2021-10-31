package model

type Comment struct {
	ID        string `json:"id"`
	PostID    string `json:"postID"`
	Comment   string `json:"comment"`
	CreatedBy string `json:"createdBy"`
	Event     string `json:"event"`
}
