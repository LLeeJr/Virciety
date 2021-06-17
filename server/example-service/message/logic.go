package message

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
)

type service struct {
	repo Repository
	logger log.Logger
}

func (s service) CreateMessage(ctx context.Context, msg string, from string, to string) (string, error) {
	logger := log.With(s.logger, "method", "CreateMessage")

	uuid, _ := uuid.NewV4()
	id := uuid.String()
	message := Message{
		ID:      id,
		Msg: msg,
		From:    from,
		To:      to,
	}

	if err := s.repo.CreateMessage(ctx, message); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("create message", id)

	return "Success", nil
}

func (s service) GetMessage(ctx context.Context, id string) (string, error) {
	logger := log.With(s.logger, "method", "GetMessage")

	message, err := s.repo.GetMessage(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("get message", id)

	return message, nil
}

func NewService(rep Repository, logger log.Logger) Service {
	return &service{
		repo:   rep,
		logger: logger,
	}
}
