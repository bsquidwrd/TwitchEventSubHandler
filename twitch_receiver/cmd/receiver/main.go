package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/cmd/receiver/routes"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/service"
)

var (
	eventsubSecret  = os.Getenv("EVENTSUBSECRET")
	eventsubWebhook = os.Getenv("EVENTSUBWEBHOOK")
)

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

	slog.Info("Receiver Starting Up...", "log_level", logLevel)

	if eventsubSecret == "" || eventsubWebhook == "" {
		panic("Verify EVENTSUBSECRET and EVENTSUBWEBHOOK are set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbServices := service.NewReceiverService()
	defer dbServices.Cleanup()

	http.HandleFunc("/", routes.HandleRoot)
	http.HandleFunc("POST /webhook", routes.HandleWebhook(dbServices))
	http.HandleFunc("GET /healthcheck", routes.HandleHealthCheck(dbServices))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Running server on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
