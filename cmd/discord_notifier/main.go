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
	discordnotifierhandlers "github.com/bsquidwrd/TwitchEventSubHandler/internal/discord_notifier_handlers"
	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "discord_notifier"

func main() {
	loggerOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOptions))
	slog.SetDefault(logger)

	slog.Info("Discord Notifier Starting Up...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbServices := database.NewDiscordNotifierService()
	defer dbServices.Cleanup()

	topics := []string{
		"stream.online",
		"stream.offline",
		"channel.update",
		"user.update",
	}

	dbServices.Queue.StartConsuming(queueName, topics, func(msg amqp.Delivery) {
		discordnotifierhandlers.ProcessMessage(dbServices, msg)
	})

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
