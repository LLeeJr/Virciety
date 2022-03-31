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
	"post-service/database"
	"post-service/graph/generated"
	"post-service/graph/model"
	"post-service/graph/resolvers"
	messagequeue "post-service/message-queue"
	"time"
)

const defaultPort = "8083"
const defaultRabbitMQUrl = "amqp://guest:guest@localhost:5672/"

func main() {
	// get rabbitmq url from envs
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

	// create repo for connection to db
	repo, err := database.NewRepo()
	messagequeue.FailOnError(err, "Failed to connect to DB")

	responses := map[string]chan []*model.Comment{}
	userResponses := map[string]chan map[string]string{}

	// create producer queue
	producerQueue, _ := messagequeue.NewPublisher()
	go producerQueue.InitPublisher(ch)

	// create consumer queue
	consumerQueue, _ := messagequeue.NewConsumer(repo, responses, userResponses)
	go consumerQueue.InitConsumer(ch)

	// graphql init
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewResolver(repo, producerQueue, responses, userResponses)}))

	// websocket config
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

	// cors config
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

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, r))
}
