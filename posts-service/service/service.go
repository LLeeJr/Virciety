package service

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"time"
)

type PostService interface {
	CreatePost(ctx context.Context, username, description, data string) (Post, error)
	GetPosts(ctx context.Context) ([]Post, error)
	RemovePost(ctx context.Context, id string) (bool, error)
	EditPost(ctx context.Context, id, newDescription string) (bool, error)
	LikedPost(ctx context.Context, id, username string) (bool, error)
}

type Service struct {
	repo   Repository
	logger log.Logger
}

func NewService(rep Repository, logger log.Logger) PostService {
	return &Service{
		repo:   rep,
		logger: logger,
	}
}

func (s Service) CreatePost(_ context.Context, username, description, data string) (Post, error) {
	logger := log.With(s.logger, "method", "CreatePost")

	if username == "" || description == "" || data == "" {
		err := errors.New("empty fields in request")
		level.Error(logger).Log("err empty fields ", err)
		return Post{}, err
	}

	post := Post{
		ID:          time.Now().String() + username,
		Description: description,
		Data:        data,
		LikedBy:     make([]string, 0),
		Comments:    make([]string, 0),
	}

	if err := s.repo.CreatePost(post); err != nil {
		level.Error(logger).Log("err from repo is ", err)
		return post, err
	}

	return post, nil
}

func (s Service) GetPosts(_ context.Context) ([]Post, error) {
	logger := log.With(s.logger, "method", "GetPosts")

	posts, err := s.repo.GetPosts()
	if err != nil {
		level.Error(logger).Log("err in getting posts ", err)
		return nil, err
	}

	return posts, nil
}

func (s Service) RemovePost(_ context.Context, id string) (bool, error) {
	logger := log.With(s.logger, "method", "RemovePost")

	removed, err := s.repo.RemovePost(id)
	if err != nil {
		level.Error(logger).Log("err ", err)
		return removed, err
	}

	return removed, nil
}

func (s Service) EditPost(_ context.Context, id, newDescription string) (bool, error) {
	logger := log.With(s.logger, "method", "EditPost")

	if id == "" || newDescription == "" {
		err := errors.New("empty fields in request")
		level.Error(logger).Log("err empty fields ", err)
		return false, err
	}

	updated, err := s.repo.EditPost(id, newDescription)
	if err != nil {
		level.Error(logger).Log("err ", err)
		return updated, err
	}

	return updated, nil
}

func (s Service) LikedPost(_ context.Context, id, username string) (bool, error) {
	logger := log.With(s.logger, "method", "LikedPost")

	if id == "" || username == "" {
		err := errors.New("empty fields in request")
		level.Error(logger).Log("err empty field ", err)
		return false, err
	}

	liked, err := s.repo.LikedPost(id, username)
	if err != nil {
		level.Error(logger).Log("err ", err)
		return liked, err
	}

	return liked, err
}