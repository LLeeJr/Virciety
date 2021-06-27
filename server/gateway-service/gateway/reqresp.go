package gateway

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	CreateDMRequest struct {
		Id    string `json:"id"`
		Msg string `json:"msg"`
	}
	CreateDMResponse struct {
		Ok string `json:"ok"`
	}

	GetDMRequest struct {
		Id string `json:"id"`
	}
	GetDMResponse struct {
		Msg string `json:"msg"`
	}
)


func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateDMReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateDMRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetDMReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetDMRequest
	vars := mux.Vars(r)

	req = GetDMRequest{
		Id: vars["id"],
	}
	return req, nil
}