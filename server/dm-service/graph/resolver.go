package graph

import (
	"dm-service/database"
	"dm-service/graph/model"
	"dm-service/queue"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func NewResolver(repo database.Repository, publisher queue.Publisher) *Resolver {
	return &Resolver{
		repo:   repo,
		publisher: publisher,
		dmChan: make(chan *model.Dm),
	}
}

type Resolver struct {
	repo      database.Repository
	publisher queue.Publisher
	dmChan    chan *model.Dm
}