package message

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateMessage endpoint.Endpoint
	GetMessage endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateMessage: makeCreateMessageEndpoint(s),
		GetMessage: makeGetMessageEndpoint(s),
	}
}

func makeCreateMessageEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateMessageRequest)
		ok, err := s.CreateMessage(ctx, req.Id, req.Msg)
		return CreateMessageResponse{Ok: ok}, err
	}
}

func makeGetMessageEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetMessageRequest)
		message, err := s.GetMessage(ctx, req.Id)
		return GetMessageResponse{Msg: message}, err
	}
}