package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func HandleNotification(notificationType string, rawBody *[]byte) {

	switch notificationType {
	case "stream.online":
		var notification models.StreamUpEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processStreamUp(notification)

	case "stream.offline":
		var notification models.StreamDownEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processStreamDown(notification)

	case "channel.update":
		var notification models.ChannelUpdateEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processChannelUpdate(notification)

	case "user.authorization.grant":
		var notification models.AuthorizationRevokeEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processAuthorizationGrant(notification)

	case "user.authorization.revoke":
		var notification models.AuthorizationRevokeEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processAuthorizationRevoke(notification)

	case "user.update":
		var notification models.UserUpdateEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processUserUpdate(notification)

	default:
		return
	}
}
