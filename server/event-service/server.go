package main

import (
	"event-service/database"
	"event-service/graph/generated"
	"event-service/graph/resolvers"
	messagequeue "event-service/message-queue"
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
	_ "net/http"
	"os"
	"time"
)

const defaultPort = "8086"

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
	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewResolver(repo, producerQueue)}))

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
