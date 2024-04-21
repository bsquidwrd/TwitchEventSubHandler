package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processAuthorizationRevoke(dbServices *database.Service, notification *models.AuthorizationRevokeEventMessage) {
	slog.Info("User revoked authorization", "userid", notification.Event.UserID)
}
