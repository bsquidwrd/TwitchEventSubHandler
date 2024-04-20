package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processUserUpdate(notification models.UserUpdateEventMessage) {
	slog.Info("User was updated", "username", notification.Event.UserName)
}
