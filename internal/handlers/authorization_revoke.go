package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/twitch"
)

func processAuthorizationRevoke(dbServices *database.Service, notification *models.AuthorizationRevokeEvent) {
	slog.Info("User revoked authorization", "userid", notification.UserID)
	defer dbServices.Queue.Publish("user.authorization.revoke", notification)

	parameters := &url.Values{}
	parameters.Add("user_id", notification.UserID)
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

	go func() {
		_, err := dbServices.Database.Exec(context.Background(), `
		delete from public.twitch_user where id=$1;
		`,
			notification.UserID,
		)

		if err != nil {
			slog.Warn("Error processing user.authorization.revoke for DB call", "id", notification.UserID)
		}
	}()
}
