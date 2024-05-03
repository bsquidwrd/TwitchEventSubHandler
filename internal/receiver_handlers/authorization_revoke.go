package receiver_handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/twitch"
	models "github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
)

func processAuthorizationRevoke(dbServices *database.ReceiverService, notification *models.AuthorizationRevokeEvent) {
	slog.Debug("User revoked authorization", "userid", notification.UserID)

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

	_, err = dbServices.Database.Exec(context.Background(), `
		delete from public.twitch_user where id=$1;
		`,
		notification.UserID,
	)

	if err != nil {
		slog.Warn("Error processing user.authorization.revoke for DB call", "user_id", notification.UserID)
		return
	}

	dbServices.Queue.Publish("user.authorization.revoke", notification)
}
