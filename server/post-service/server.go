package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"posts-service/database"
	"posts-service/graph/generated"
	"posts-service/graph/model"
	"posts-service/graph/resolvers"
	messagequeue "posts-service/message-queue"
	"time"
)

const defaultPort = "8083"
const defaultRabbitMQUrl = "amqp://guest:guest@localhost:5672/"

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

	responses := map[string]chan []*model.Comment{}

	producerQueue, _ := messagequeue.NewPublisher()
	go producerQueue.InitPublisher(ch)

	consumerQueue, _ := messagequeue.NewConsumer(repo, responses)
	go consumerQueue.InitConsumer(ch)

	// graphql init
	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewResolver(repo, producerQueue)}))

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewResolver(repo, producerQueue, responses)}))

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
