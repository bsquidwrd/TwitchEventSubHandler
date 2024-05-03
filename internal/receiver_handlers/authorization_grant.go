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

func processAuthorizationGrant(dbServices *database.ReceiverService, notification *models.AuthorizationGrantEvent) {
	slog.Info("User granted authorization", "userid", notification.UserID)

	// Get User information
	parameters := &url.Values{}
	parameters.Add("id", notification.UserID)
	_, twitchApiUser, err := twitch.CallApi(dbServices, http.MethodGet, "users", "", parameters)
	if err != nil {
		slog.Warn("Could not retrieve user from Twitch API", "error", err)
		return
	}

	var twitchUsers models.UserData
	json.Unmarshal(twitchApiUser, &twitchUsers)

	if len(twitchUsers.Data) > 1 {
		slog.Warn("Multiple users returned when requesting from Twitch API", "user_id", notification.UserID, "body", string(twitchApiUser))
		return
	} else if len(twitchUsers.Data) < 1 {
		slog.Warn("No users returned when requesting from Twitch API", "user_id", notification.UserID, "body", string(twitchApiUser))
		return
	}

	twitchUser := twitchUsers.Data[0]

	// Get Channel information
	parameters = &url.Values{}
	parameters.Add("broadcaster_id", notification.UserID)
	_, twitchApiChannel, err := twitch.CallApi(dbServices, http.MethodGet, "channels", "", parameters)
	if err != nil {
		slog.Warn("Could not retrieve user from Twitch API", "error", err)
		return
	}

	var twitchChannels models.ChannelData
	json.Unmarshal(twitchApiChannel, &twitchChannels)

	if len(twitchUsers.Data) > 1 {
		slog.Warn("Multiple channels returned when requesting from Twitch API", "user_id", notification.UserID, "body", string(twitchApiChannel))
		return
	} else if len(twitchUsers.Data) < 1 {
		slog.Warn("No channels returned when requesting from Twitch API", "user_id", notification.UserID, "body", string(twitchApiChannel))
	}

	twitchChannel := twitchChannels.Data[0]

	// Save info to database
	_, err = dbServices.Database.Exec(context.Background(), `
		insert into public.twitch_user (id,"name",login,avatar_url,email,description,title,"language",category_id,category_name)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		on conflict (id) do update
		set "name"=$2,login=$3,avatar_url=$4,email=$5,description=$6,title=$7,"language"=$8,category_id=$9,category_name=$10
		`,
		twitchUser.ID,
		twitchUser.DisplayName,
		twitchUser.Login,
		twitchUser.ProfileImageUrl,
		twitchUser.Email,
		twitchUser.Description,
		twitchChannel.Title,
		twitchChannel.BroadcasterLanguage,
		twitchChannel.GameID,
		twitchChannel.GameName,
	)

	if err != nil {
		slog.Warn("Error processing user.authorization.grant for DB call", "user_id", notification.UserID)
		return
	}

	// Subscribe to other events of interest
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
		subType := subscription.Type
		go func() {
			_, response, err := twitch.CallApi(dbServices, http.MethodPost, "eventsub/subscriptions", body, nil)
			if err != nil {
				slog.Warn("Could not subscribe to event for user", "subscription_type", subType, "error", err, "response", string(response))
			}
		}()
	}

	dbServices.Queue.Publish("user.authorization.grant", notification)
}
