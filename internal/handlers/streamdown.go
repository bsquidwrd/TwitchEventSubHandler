package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processStreamDown(dbServices *database.Service, notification *models.StreamDownEventMessage) {
	slog.Info("Channel went offline", "username", notification.Event.BroadcasterUserName)
	defer dbServices.Queue.Publish("stream.offline", notification)
}
