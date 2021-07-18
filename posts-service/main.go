package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
	_ "net/http"
	"os"
	"posts-service/database"
	"posts-service/graph/generated"
	"posts-service/graph/resolvers"
	messagequeue "posts-service/message-queue"
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

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewResolver(repo, producerQueue)}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
