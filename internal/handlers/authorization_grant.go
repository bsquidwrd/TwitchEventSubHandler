package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/twitch"
)

func processAuthorizationGrant(dbServices *database.Service, notification *models.AuthorizationGrantEventMessage) {
	slog.Info("User granted authorization", "userid", notification.Event.UserID)

	eventsubSecret := os.Getenv("EVENTSUBSECRET")
	eventsubWebhook := os.Getenv("EVENTSUBWEBHOOK")
	bodies := []string{
		fmt.Sprintf(`{"type":"user.update","version":"1","condition":{"user_id":"%s"},"transport":{"method":"webhook","callback":"%s","secret":"%s"}}`, notification.Event.UserID, eventsubWebhook, eventsubSecret),
		fmt.Sprintf(`{"type":"channel.update","version":"2","condition":{"broadcaster_user_id":"%s"},"transport":{"method":"webhook","callback":"%s","secret":"%s"}}`, notification.Event.UserID, eventsubWebhook, eventsubSecret),
		fmt.Sprintf(`{"type":"stream.online","version":"1","condition":{"broadcaster_user_id":"%s"},"transport":{"method":"webhook","callback":"%s","secret":"%s"}}`, notification.Event.UserID, eventsubWebhook, eventsubSecret),
		fmt.Sprintf(`{"type":"stream.offline","version":"1","condition":{"broadcaster_user_id":"%s"},"transport":{"method":"webhook","callback":"%s","secret":"%s"}}`, notification.Event.UserID, eventsubWebhook, eventsubSecret),
	}

	for _, body := range bodies {
		go twitch.CallApi(dbServices, http.MethodPost, "eventsub/subscriptions", body, nil)
	}
}
