package gateway

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-kit/kit/log"
)

var RepoErr = errors.New("unable to handle repo request")

type repo struct {
	db *sql.DB
	logger log.Logger
}

func (r repo) CreateDMEvent(ctx context.Context, event DMEvent) error {


	// post via graphql to database
	return nil
}

func (r repo) GetDM(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", RepoErr
	}

	// get from db via graphql
	result := ""

	return result, nil
}

func NewRepo(db *sql.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}
