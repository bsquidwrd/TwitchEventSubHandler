package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processChannelUpdate(dbServices *database.Services, notification models.ChannelUpdateEventMessage) {
	slog.Info("Channel was updated", "username", notification.Event.BroadcasterUserName)
}
