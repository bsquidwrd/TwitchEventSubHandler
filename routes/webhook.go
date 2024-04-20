package routes

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/handlers"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/utils"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("SECRET")

	if secret == "" {
		slog.Error("Secret could not be found")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Warn("Error reading body from request", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messageID := r.Header.Get("Twitch-Eventsub-Message-Id")
	messageSignature := r.Header.Get("Twitch-Eventsub-Message-Signature")[7:]
	messageTimestamp := r.Header.Get("Twitch-Eventsub-Message-Timestamp")

	if !utils.ValidateSignature([]byte(secret), messageID, messageTimestamp, rawBody, messageSignature) {
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
		if eventsubMessage.Subscription.Cost > 0 &&
			eventsubMessage.Subscription.Type != "user.authorization.grant" &&
			eventsubMessage.Subscription.Type != "user.authorization.revoke" {
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

	case "notification":
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, "Oh that's lit!")

		go handlers.HandleNotification(r, &rawBody)

	case "revocation":
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, "Such is life")

		go handlers.HandleRevocation(r, &rawBody)

	default:
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Womp womp")
	}
}
