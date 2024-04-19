package helpers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bsquidwrd/twitcheventsub-receiver/models"
)

func HandleNotification(r *http.Request, rawBody *[]byte) {
	notificationType := r.Header.Get("Twitch-Eventsub-Subscription-Type")

	switch notificationType {
	case "stream.online":
		var notification models.StreamUpEvent
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		slog.Info("User went live", "username", notification.Event.BroadcasterUserName)

	default:
		return
	}
}
