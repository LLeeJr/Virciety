package message

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/message").Handler(httptransport.NewServer(
		endpoints.CreateMessage,
		decodeCreateMessageReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/message/{id}").Handler(httptransport.NewServer(
		endpoints.GetMessage,
		decodeGetMessageReq,
		encodeResponse,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}