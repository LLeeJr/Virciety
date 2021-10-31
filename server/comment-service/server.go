package main

import (
	"comment-service/database"
	"comment-service/graph/generated"
	"comment-service/graph/resolvers"
	messagequeue "comment-service/message-queue"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	_ "net/http"
	"os"
)

const defaultPort = "8084"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	repo, _ := database.NewRepo()

	// rabbitmq connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	messagequeue.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	messagequeue.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	producerQueue, _ := messagequeue.NewPublisher()
	go producerQueue.InitPublisher(ch)

	consumerQueue, _ := messagequeue.NewConsumer(repo)
	go consumerQueue.InitConsumer(ch)

	// graphql init
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewResolver(repo, producerQueue)}))

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
