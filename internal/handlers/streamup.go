package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processStreamUp(notification models.StreamUpEventMessage) {
	slog.Info("Channel went live", "username", notification.Event.BroadcasterUserName)
}
