package graph

import (
	"notifs-service/database"
	"notifs-service/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func NewResolver(repo database.Repository) *Resolver {
	return &Resolver{
		repo:       repo,
		notifsChan: make(chan *model.Notif),
	}
}

type Resolver struct {
	repo       database.Repository
	notifs     []*model.Notif
	notifsChan chan *model.Notif
}
