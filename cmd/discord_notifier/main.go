package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/cmd/discord_notifier/routes"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "discord_notifier"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Discord Notifier Starting Up...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = "8081"

	dbServices := database.NewDiscordNotifierService()
	defer dbServices.Cleanup()

	topics := []string{
		"stream.online",
		"stream.offline",
		"channel.update",
		"user.update",
	}

	dbServices.Queue.StartConsuming(queueName, topics, testFunc)

	http.HandleFunc("/", routes.HandleRoot)
	http.HandleFunc("GET /healthcheck", routes.HandleHealthCheck(dbServices))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Running server on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func testFunc(msg amqp.Delivery) {
	slog.Info("Got message", "topic", msg.RoutingKey)
	<-time.After(5 * time.Second)
}
