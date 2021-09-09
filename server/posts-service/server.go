package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"log"
	"net/http"
	_ "net/http"
	"os"
	"posts-service/database"
	"posts-service/graph/generated"
	"posts-service/graph/resolvers"
	messagequeue "posts-service/message-queue"
	"time"
)

const defaultPort = "8083"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := database.GetDBConn()

	repo, _ := database.NewRepo(db)

	producerQueue, _ := messagequeue.NewPublisher()
	go producerQueue.InitPublisher()

	consumerQueue, _ := messagequeue.NewConsumer(repo)
	go consumerQueue.InitConsumer()

	//srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewResolver(repo, producerQueue)}))

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewResolver(repo, producerQueue)}))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.Use(extension.Introspection{})

	r := mux.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	r.Handle("/query", srv)
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
