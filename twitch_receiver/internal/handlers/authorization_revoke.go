package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/api"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/service"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/pkg/models"
)

func processAuthorizationRevoke(dbServices *service.ReceiverService, notification *models.AuthorizationRevokeEvent) {
	slog.Debug("User revoked authorization", "userid", notification.UserID)

	parameters := &url.Values{}
	parameters.Add("user_id", notification.UserID)
	_, rawBody, err := api.CallApi(dbServices, http.MethodGet, "eventsub/subscriptions", "", parameters)
	if err != nil {
		return
	}

	var subscriptions *models.EventsubSubscriptionList
	json.Unmarshal(rawBody, &subscriptions)

	for _, subscription := range subscriptions.Data {
		if subscription.Transport.Callback != os.Getenv("EVENTSUBWEBHOOK") {
			continue
		}
		go api.DeleteSubscription(dbServices, subscription.ID)
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
