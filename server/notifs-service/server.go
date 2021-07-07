package main

import (
	"log"
	"net/http"
	"notifs-service/database"
	"notifs-service/graph"
	"notifs-service/graph/generated"
	"notifs-service/queue"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
)

const defaultPort = "8082"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := database.Connect()
	repo, err := database.NewRepo(db)
	if err != nil {
		log.Fatal("err", err)
	}

	publisher, _ := queue.NewPublisher()
	go publisher.InitPublisher()

	consumer, _ := queue.NewConsumer(repo)
	go consumer.InitConsumer()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(repo, publisher)}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
