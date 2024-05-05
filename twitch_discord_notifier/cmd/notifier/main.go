package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_discord_notifier/cmd/notifier/routes"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_discord_notifier/internal/handlers"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_discord_notifier/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "discord_notifier"

func main() {
	debugEnabled := os.Getenv("DEBUG")
	logLevel := slog.LevelInfo

	if debugEnabled != "" {
		logLevel = slog.LevelDebug
	}
	loggerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOptions))
	slog.SetDefault(logger)

	slog.Info("Discord Notifier Starting Up...", "log_level", logLevel)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbServices := service.NewDiscordNotifierService()
	defer dbServices.Cleanup()

	topics := []string{
		"stream.online",
		"stream.offline",
		"channel.update",
		"user.update",
	}

	dbServices.Queue.StartConsuming(queueName, topics, func(msg amqp.Delivery) {
		handlers.ProcessMessage(dbServices, msg)
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
