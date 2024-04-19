package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/bsquidwrd/TwitchEventSubHandler/routes"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	secret := os.Getenv("SECRET")
	port := os.Getenv("PORT")

	if secret == "" || port == "" {
		log.Fatal("You must specify SECRET and PORT environment variables")
	}

	fmt.Println("Starting up!")

	http.HandleFunc("POST /webhook", routes.HandleWebhook)
	http.HandleFunc("GET /healthcheck", routes.HandleHealthCheck)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
