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
		var notification models.StreamUpEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		slog.Info("Channel went live", "username", notification.Event.BroadcasterUserName)

	case "stream.offline":
		var notification models.StreamDownEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		slog.Info("Channel went offline", "username", notification.Event.BroadcasterUserName)

	case "channel.update":
		var notification models.ChannelUpdateEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		slog.Info("Channel was updated", "username", notification.Event.BroadcasterUserName)

	case "user.update":
		var notification models.UserUpdateEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		slog.Info("User was updated", "username", notification.Event.UserName)

	default:
		return
	}
}
