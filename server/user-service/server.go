package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"user-service/database"
	"user-service/graph"
	"user-service/graph/generated"
	"user-service/queue"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8085"
const defaultRabbitMQUrl = "amqp://guest:guest@localhost:5672"

// main function for starting the dm-microservice
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = defaultRabbitMQUrl
	}

	// rabbitmq connection
	conn, err := amqp.Dial(rabbitmqURL)
	queue.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	queue.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	repo, err := database.NewRepo()
	if err != nil {
		log.Fatal("err", err)
	}

	publisher, _ := queue.NewPublisher()
	go publisher.InitPublisher(ch)

	consumer, _ := queue.NewConsumer(repo)
	go consumer.InitConsumer(ch)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(repo, publisher)}))

	r := mux.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
