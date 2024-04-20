package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processChannelUpdate(notification models.ChannelUpdateEventMessage) {
	slog.Info("Channel was updated", "username", notification.Event.BroadcasterUserName)
}
