package graph

import (
	"dm-service/database"
	"dm-service/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func NewResolver(repo database.Repository) *Resolver {
	return &Resolver{
		repo: repo,
		dmsChan: make(chan *model.Dm),
	}

}

type Resolver struct {
	repo 	   database.Repository
	dms        []*model.Dm
	dmsChan chan *model.Dm
}