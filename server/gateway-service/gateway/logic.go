package gateway

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"strings"
	"time"
)

type service struct {
	repository Repository
	logger log.Logger
}

func (s service) CreateDMEvent(ctx context.Context, id string, msg string) (string, error) {
	logger := log.With(s.logger, "method", "CreateDM")

	data := strings.Split(id, "__")

	dm := DM{
		ID:  id,
		Msg: msg,
	}

	dmEvent := DMEvent{
		EventTime: time.Now().String(),
		From:      data[0],
		To:        data[1],
		Time:      data[2],
		Payload:   dm,
	}

	logger.Log("event", dmEvent)

	// todo: store event, enqueue createDMRequest

	err := s.repository.CreateDMEvent(ctx, dmEvent)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("create dm", id)
	return "Success", err
}

func (s service) GetDMEvent(ctx context.Context, id string) (string, error) {
	logger := log.With(s.logger, "method", "GetDM")

	msg, err := s.repository.GetDM(ctx, id)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("get dm", id)

	return msg, nil
}

func NewService(rep Repository, logger log.Logger) Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}
