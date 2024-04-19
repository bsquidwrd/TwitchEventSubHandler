package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bsquidwrd/twitcheventsub-receiver/helpers"
	"github.com/bsquidwrd/twitcheventsub-receiver/models"
)

var secret string = "0123456789"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	fmt.Println("Starting up!")

	http.HandleFunc("/webhook", handleWebhook)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messageID := r.Header.Get("Twitch-Eventsub-Message-Id")
	messageSignature := r.Header.Get("Twitch-Eventsub-Message-Signature")[7:]
	messageTimestamp := r.Header.Get("Twitch-Eventsub-Message-Timestamp")

	if !helpers.ValidateSignature([]byte(secret), messageID, messageTimestamp, rawBody, messageSignature) {
		slog.Warn("Invalid request received", "endpoint", html.EscapeString(r.URL.Path))
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Womp womp")
		return
	}

	slog.Info(
		"Successful request received",
		"endpoint", html.EscapeString(r.URL.Path),
		"type", r.Header.Get("Twitch-Eventsub-Message-Type"),
		"subscription", r.Header.Get("Twitch-Eventsub-Subscription-Type"),
	)

	w.Header().Add("Content-Type", "text/plain")
	switch r.Header.Get("Twitch-Eventsub-Message-Type") {

	case "webhook_callback_verification":
		var challenge models.EventsubSubscriptionVerification
		err = json.Unmarshal(rawBody, &challenge)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, challenge.Challenge)

	case "revocation":
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, "Such is life")

	case "notification":
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, "Oh that's lit!")

		go helpers.HandleNotification(r, &rawBody)

	default:
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Womp womp")
	}
}
