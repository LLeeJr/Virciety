package main

import (
	"dm-service/database"
	"dm-service/graph"
	"dm-service/graph/generated"
	"dm-service/queue"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
)

const defaultPort = "8081"

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
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
