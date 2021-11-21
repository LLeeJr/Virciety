// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type User struct {
	ID        string   `json:"id"`
	Username  string   `json:"username"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Follows   []string `json:"follows"`
}

type UserData struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}