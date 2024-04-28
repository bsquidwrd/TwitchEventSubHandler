package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/twitch"
	"github.com/bsquidwrd/TwitchEventSubHandler/pkg/models"
)

func processAuthorizationGrant(dbServices *database.ReceiverService, notification *models.AuthorizationGrantEvent) {
	slog.Info("User granted authorization", "userid", notification.UserID)
	defer dbServices.Queue.Publish("user.authorization.grant", notification)

	eventsubSecret := os.Getenv("EVENTSUBSECRET")
	eventsubWebhook := os.Getenv("EVENTSUBWEBHOOK")
	subscriptions := []models.EventsubSubscription{
		{
			Type:    "user.update",
			Version: "1",
			Condition: models.EventsubCondition{
				UserID: notification.UserID,
			},
			Transport: models.EventsubTransport{
				Method:   "webhook",
				Callback: eventsubWebhook,
				Secret:   eventsubSecret,
			},
		},
		{
			Type:    "channel.update",
			Version: "2",
			Condition: models.EventsubCondition{
				BroadcasterUserID: notification.UserID,
			},
			Transport: models.EventsubTransport{
				Method:   "webhook",
				Callback: eventsubWebhook,
				Secret:   eventsubSecret,
			},
		},
		{
			Type:    "stream.online",
			Version: "1",
			Condition: models.EventsubCondition{
				BroadcasterUserID: notification.UserID,
			},
			Transport: models.EventsubTransport{
				Method:   "webhook",
				Callback: eventsubWebhook,
				Secret:   eventsubSecret,
			},
		},
		{
			Type:    "stream.offline",
			Version: "1",
			Condition: models.EventsubCondition{
				BroadcasterUserID: notification.UserID,
			},
			Transport: models.EventsubTransport{
				Method:   "webhook",
				Callback: eventsubWebhook,
				Secret:   eventsubSecret,
			},
		},
	}

	for _, subscription := range subscriptions {
		bodyBytes, _ := json.Marshal(subscription)
		body := string(bodyBytes)
		go twitch.CallApi(dbServices, http.MethodPost, "eventsub/subscriptions", body, nil)
	}

	go func() {
		_, err := dbServices.Database.Exec(context.Background(), `
		insert into public.twitch_user (id,"name",login)
		values($1,$2,$3)
		on conflict (id) do update
		set "name"=$2,login=$3;
		`,
			notification.UserID,
			notification.UserName,
			notification.UserLogin,
		)

		if err != nil {
			slog.Warn("Error processing user.authorization.grant for DB call", "userid", notification.UserID)
		}
	}()
}
