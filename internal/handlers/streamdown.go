package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processStreamDown(notification models.StreamDownEventMessage) {
	slog.Info("Channel went offline", "username", notification.Event.BroadcasterUserName)
}
