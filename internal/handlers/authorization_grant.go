package handlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processAuthorizationGrant(notification models.AuthorizationRevokeEventMessage) {
	slog.Info("User granted authorization", "userid", notification.Event.UserID)
}
