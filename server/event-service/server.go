package main

import (
	"event-service/database"
	"event-service/graph/generated"
	"event-service/graph/resolvers"
	messagequeue "event-service/message-queue"
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

const defaultPort = "8086"
const defaultRabbitMQUrl = "amqp://guest:guest@localhost:5672"

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
	messagequeue.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	messagequeue.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	repo, err := database.NewRepo()
	messagequeue.FailOnError(err, "Failed to connect to DB")

	producerQueue, _ := messagequeue.NewPublisher()
	go producerQueue.InitPublisher(ch)

	consumerQueue, _ := messagequeue.NewConsumer(repo)
	go consumerQueue.InitConsumer(ch)

	// graphql init
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewResolver(repo, producerQueue)}))

	r := mux.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	r.Handle("/query", srv)
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
