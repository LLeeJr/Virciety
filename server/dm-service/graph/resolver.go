package graph

import (
	"context"
	"dm-service/database"
	"dm-service/graph/generated"
	"dm-service/queue"
	"github.com/99designs/gqlgen/graphql"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ckey string

func NewResolver(repo database.Repository, publisher queue.Publisher) *Resolver {
	return &Resolver{
		Rooms: map[string]*Chatroom{},
		repo:   repo,
		publisher: publisher,
	}
}

func New(repo database.Repository, publisher queue.Publisher) generated.Config {
	return generated.Config{
		Resolvers: &Resolver{
			Rooms: map[string]*Chatroom{},
			repo:   repo,
			publisher: publisher,
		},
		Directives: generated.DirectiveRoot{
			User: func(ctx context.Context, obj interface{}, next graphql.Resolver, username string) (res interface{}, err error) {
				return next(context.WithValue(ctx, ckey("username"), username))
			},
		},
	}
}

type Resolver struct {
	Rooms map[string]*Chatroom
	mu sync.Mutex
	repo      database.Repository
	publisher queue.Publisher
}

func getUsername(ctx context.Context) string {
	if username, ok := ctx.Value(ckey("username")).(string); ok {
		return username
	}
	return ""
}