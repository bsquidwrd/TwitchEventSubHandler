package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processAuthorizationRevoke(notification models.AuthorizationRevokeEventMessage) {
	slog.Info("User revoked authorization", "userid", notification.Event.UserID)
}
