package helpers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bsquidwrd/twitcheventsub-receiver/models"
)

func HandleRevocation(r *http.Request, rawBody *[]byte) {
	var revocation models.AuthorizationRevokeEventMessage
	err := json.Unmarshal(*rawBody, &revocation)
	if err != nil {
		slog.Error("Could not unmarshal body", err)
		return
	}

	slog.Info("User revoked access to application", "username", revocation.Event.UserName)
}
