// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateDmRequest struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

type Dm struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

type GetChatRequest struct {
	User1 string `json:"user1"`
	User2 string `json:"user2"`
}
