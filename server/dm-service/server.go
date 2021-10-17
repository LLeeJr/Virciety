package main

import (
	"dm-service/database"
	"dm-service/graph"
	"dm-service/graph/generated"
	"dm-service/queue"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultPort = "8081"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	repo, err := database.NewRepo()
	if err != nil {
		log.Fatal("err", err)
	}

	publisher, _ := queue.NewPublisher()
	go publisher.InitPublisher()

	consumer, _ := queue.NewConsumer(repo)
	go consumer.InitConsumer()

	//srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(repo, publisher)}))
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(repo, publisher)}))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		Upgrader:              websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	r := mux.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedMethods: []string{"GET","POST", "OPTIONS"},
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedHeaders: []string{"Authorization","Content-Type","Bearer","Bearer ","content-type","Origin","Accept"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
