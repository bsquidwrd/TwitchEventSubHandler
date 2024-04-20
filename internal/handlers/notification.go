package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func HandleNotification(dbServices *database.Services, notificationType string, rawBody *[]byte) {

	switch notificationType {
	case "stream.online":
		var notification models.StreamUpEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processStreamUp(dbServices, notification)

	case "stream.offline":
		var notification models.StreamDownEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processStreamDown(dbServices, notification)

	case "channel.update":
		var notification models.ChannelUpdateEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processChannelUpdate(dbServices, notification)

	case "user.authorization.grant":
		var notification models.AuthorizationGrantEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processAuthorizationGrant(dbServices, notification)

	case "user.authorization.revoke":
		var notification models.AuthorizationRevokeEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processAuthorizationRevoke(dbServices, notification)

	case "user.update":
		var notification models.UserUpdateEventMessage
		err := json.Unmarshal(*rawBody, &notification)
		if err != nil {
			slog.Error("Could not unmarshal body", err)
			return
		}

		processUserUpdate(dbServices, notification)

	default:
		return
	}
}
