package gateway

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/dm").Handler(httptransport.NewServer(
		endpoints.CreateDM,
		decodeCreateDMReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/dm/{id}").Handler(httptransport.NewServer(
		endpoints.GetDM,
		decodeGetDMReq,
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