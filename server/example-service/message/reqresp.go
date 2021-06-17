package message

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	CreateMessageRequest struct {
		Msg string `json:"msg"`
		From string `json:"from"`
		To string `json:"to"`
	}
	CreateMessageResponse struct {
		Ok string `json:"ok"`
	}
	
	GetMessageRequest struct {
		Id string `json:"id"`
	}
	GetMessageResponse struct {
		Msg string `json:"msg"`
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateMessageReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateMessageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetMessageReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetMessageRequest
	vars := mux.Vars(r)

	req = GetMessageRequest{Id: vars["id"]}

	return req, nil
}