package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/twitch"
)

func processAuthorizationRevoke(dbServices *database.Service, notification *models.AuthorizationRevokeEventMessage) {
	slog.Info("User revoked authorization", "userid", notification.Event.UserID)

	parameters := &url.Values{}
	parameters.Add("user_id", notification.Event.UserID)
	_, rawBody, err := twitch.CallApi(dbServices, http.MethodGet, "eventsub/subscriptions", "", parameters)
	if err != nil {
		return
	}

	var subscriptions *models.EventsubSubscriptionList
	json.Unmarshal(rawBody, &subscriptions)

	for _, subscription := range subscriptions.Data {
		if subscription.Transport.Callback != os.Getenv("EVENTSUBWEBHOOK") {
			continue
		}
		go twitch.DeleteSubscription(dbServices, subscription.ID)
	}
}
