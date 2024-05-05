package routes

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/api"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/handlers"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/service"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/utils"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/pkg/models"
)

func HandleWebhook(dbServices *service.ReceiverService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("EVENTSUBSECRET")

		if secret == "" {
			slog.Error("EVENTSUBSECRET could not be found")
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
		messageType := r.Header.Get("Twitch-Eventsub-Message-Type")

		parsedTimestamp, err := time.Parse(time.RFC3339, messageTimestamp)
		if err != nil {
			slog.Warn("Invalid Message Timestamp detected", "timestamp", messageTimestamp, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// If a message is older than 10 minutes, don't do anything with it
		if time.Now().UTC().Sub(parsedTimestamp).Minutes() >= 10 {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "This is as outdated as flip phones.")
			return
		}

		// Ensure we haven't processed this message id already
		if dbServices.Cache.GetBool(fmt.Sprintf("twitch:message:%s", messageID)) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "I've seen this so much, it's starting to feel like deja-vu-vu.")
			return
		} else {
			dbServices.Cache.SetBool(fmt.Sprintf("twitch:message:%s", messageID), true, 10*time.Minute)
		}

		if !utils.ValidateSignature([]byte(secret), messageID, messageTimestamp, rawBody, messageSignature) {
			slog.Warn("Invalid request received", "endpoint", html.EscapeString(r.URL.Path))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		slog.Debug(
			"Successful request received",
			"endpoint", html.EscapeString(r.URL.Path),
			"type", messageType,
			"subscription", r.Header.Get("Twitch-Eventsub-Subscription-Type"),
		)

		w.Header().Add("Content-Type", "text/plain")

		// Verify we even have the right kind of base body
		var eventsubMessage models.EventsubMessage
		err = json.Unmarshal(rawBody, &eventsubMessage)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Depending on the message, handle it differently
		switch messageType {

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

		case "notification", "revocation":
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, "Oh that's lit!")

			go handlers.HandleNotification(dbServices, r.Header.Get("Twitch-Eventsub-Subscription-Type"), &rawBody)

		default:
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Womp womp")
		}

		// Perform cost check
		// If there's a cost with a revocation, that's okay
		if eventsubMessage.Subscription.Cost > 0 &&
			r.Header.Get("Twitch-Eventsub-Message-Type") != "revocation" &&
			eventsubMessage.Subscription.Type != "user.authorization.grant" &&
			eventsubMessage.Subscription.Type != "user.authorization.revoke" {
			go api.DeleteSubscription(dbServices, eventsubMessage.Subscription.ID)
		}
	}
}
