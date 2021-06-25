package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	_ "net/http"
	"os"
	"posts-service/service"
	"posts-service/service/database"
	"posts-service/service/transport"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	r := mux.NewRouter()
	db := database.GetDBConn()

	var svc service.PostService
	svc = service.Service{}
	{
		repository, err := service.NewRepo(db, logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		svc = service.NewService(repository, logger)
	}

	CreatePostHandler := httptransport.NewServer(
		transport.MakeCreatePostEndpoint(svc),
		transport.DecodeCreatePostRequest,
		transport.EncodeResponse)

	GetPostsHandler := httptransport.NewServer(
		transport.MakeGetPostsEndpoint(svc),
		transport.DecodeGetPostsRequest,
		transport.EncodeResponse)

	RemovePostHandler := httptransport.NewServer(
		transport.MakeRemovePostEndpoint(svc),
		transport.DecodeRemovePostRequest,
		transport.EncodeResponse)

	EditPostHandler := httptransport.NewServer(
		transport.MakeEditPostEndpoint(svc),
		transport.DecodeEditPostRequest,
		transport.EncodeResponse)

	LikedPostHandler := httptransport.NewServer(
		transport.MakeLikedPostEndpoint(svc),
		transport.DecodeLikedPostRequest,
		transport.EncodeResponse)

	http.Handle("/", r)
	http.Handle("/post", CreatePostHandler)
	http.Handle("/post/edit", EditPostHandler)
	http.Handle("/post/liked", LikedPostHandler)
	r.Handle("/post/getAll", GetPostsHandler).Methods("GET")
	r.Handle("/post/{id}", RemovePostHandler).Methods("DELETE")

	logger.Log("msg", "HTTP", "addr", ":8000")
	logger.Log("err", http.ListenAndServe(":8000", nil))

}
