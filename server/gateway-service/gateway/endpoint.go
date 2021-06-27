package gateway

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateDM endpoint.Endpoint
	GetDM endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateDM: makeCreateDMEndpoint(s),
		GetDM:    makeGetDMEndpoint(s),
	}
}

func makeGetDMEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetDMRequest)
		msg, err := s.GetDMEvent(ctx, req.Id)

		return GetDMResponse{Msg: msg}, err
	}
}

func makeCreateDMEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateDMRequest)
		ok, err := s.CreateDMEvent(ctx, req.Id, req.Msg)
		return CreateDMResponse{Ok: ok}, err
	}
}

