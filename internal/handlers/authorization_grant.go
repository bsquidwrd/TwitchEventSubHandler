package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processAuthorizationGrant(dbServices *database.Services, notification models.AuthorizationRevokeEventMessage) {
	slog.Info("User granted authorization", "userid", notification.Event.UserID)
}
