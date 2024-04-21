package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processStreamUp(dbServices *database.Service, notification *models.StreamUpEventMessage) {
	slog.Info("Channel went live", "username", notification.Event.BroadcasterUserName)
}
