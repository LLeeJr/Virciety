// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateNotifRequest struct {
	Receiver string `json:"receiver"`
	Text     string `json:"text"`
	Event    string `json:"event"`
}

type Notif struct {
	ID       string `json:"id"`
	Receiver string `json:"receiver"`
	Text     string `json:"text"`
	Event    string `json:"event"`
}