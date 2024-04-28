package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/bsquidwrd/TwitchEventSubHandler/cmd/receiver/routes"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Receiver Starting Up...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbServices := database.NewReceiverService()
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
