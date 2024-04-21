package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processUserUpdate(dbServices *database.Service, notification *models.UserUpdateEventMessage) {
	slog.Info("User was updated", "username", notification.Event.UserName)
}
