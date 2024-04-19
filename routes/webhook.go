package routes

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/bsquidwrd/TwitchEventSubHandler/helpers"
	"github.com/bsquidwrd/TwitchEventSubHandler/models"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var secret string = os.Getenv("SECRET")

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

	// Perform cost check
	// If there's a cost with a revocation, that's okay
	// We'll be cleaning it up anyway
	if r.Header.Get("Twitch-Eventsub-Message-Type") != "revocation" {
		var eventsubMessage models.EventsubMessage
		err = json.Unmarshal(rawBody, &eventsubMessage)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if eventsubMessage.Subscription.Cost > 0 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "That's too rich for my blood")
			return
		}
	}

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

		go helpers.HandleRevocation(r, &rawBody)

	case "notification":
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, "Oh that's lit!")

		go helpers.HandleNotification(r, &rawBody)

	default:
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Womp womp")
	}
}