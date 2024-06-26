package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	models "github.com/bsquidwrd/TwitchEventSubHandler/shared/models/eventsub"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/api"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/service"
)

func processAuthorizationGrant(dbServices *service.ReceiverService, notification *models.AuthorizationGrantEvent) {
	slog.Debug("User granted authorization", "userid", notification.UserID)

	// Get User information
	parameters := &url.Values{}
	parameters.Add("id", notification.UserID)
	_, twitchApiUser, err := api.CallApi(dbServices, http.MethodGet, "users", "", parameters)
	if err != nil {
		slog.Warn("Could not retrieve user from Twitch API", "error", err)
		return
	} else {
		slog.Debug("Successfully got data from API for users")
	}

	var twitchUsers models.UserData
	err = json.Unmarshal(twitchApiUser, &twitchUsers)
	if err != nil {
		slog.Warn("Could not unmarshal users response", "error", err)
	} else {
		slog.Debug("Successfully unmarshaled API data for users")
	}

	if len(twitchUsers.Data) > 1 {
		slog.Warn("Multiple users returned when requesting from Twitch API", "user_id", notification.UserID, "body", string(twitchApiUser))
		return
	} else if len(twitchUsers.Data) < 1 {
		slog.Warn("No users returned when requesting from Twitch API", "user_id", notification.UserID, "body", string(twitchApiUser))
		return
	} else {
		slog.Debug("Only got 1 user from Twitch API")
	}

	twitchUser := twitchUsers.Data[0]

	// Get Channel information
	parameters = &url.Values{}
	parameters.Add("broadcaster_id", notification.UserID)
	_, twitchApiChannel, err := api.CallApi(dbServices, http.MethodGet, "channels", "", parameters)
	if err != nil {
		slog.Warn("Could not retrieve user from Twitch API", "error", err)
		return
	} else {
		slog.Debug("Successfully got data from API for channels")
	}

	var twitchChannels models.ChannelData
	err = json.Unmarshal(twitchApiChannel, &twitchChannels)
	if err != nil {
		slog.Warn("Could not unmarshal users response", "error", err)
	} else {
		slog.Debug("Successfully unmarshaled data from API for channels")
	}

	if len(twitchChannels.Data) > 1 {
		slog.Warn("Multiple channels returned when requesting from Twitch API", "user_id", notification.UserID, "body", string(twitchApiChannel))
		return
	} else if len(twitchChannels.Data) < 1 {
		slog.Warn("No channels returned when requesting from Twitch API", "user_id", notification.UserID, "body", string(twitchApiChannel))
	} else {
		slog.Debug("Only got one channel from API response")
	}

	twitchChannel := twitchChannels.Data[0]

	// Save info to database
	_, err = dbServices.Database.Exec(context.Background(), `
		insert into public.twitch_user (id,"name",login,avatar_url,description,title,"language",category_id,category_name)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9)
		on conflict (id) do update
		set "name"=$2,login=$3,avatar_url=$4,description=$5,title=$6,"language"=$7,category_id=$8,category_name=$9
		`,
		twitchUser.ID,
		twitchUser.DisplayName,
		twitchUser.Login,
		twitchUser.ProfileImageUrl,
		twitchUser.Description,
		twitchChannel.Title,
		twitchChannel.BroadcasterLanguage,
		twitchChannel.GameID,
		twitchChannel.GameName,
	)

	if err != nil {
		slog.Warn("Error processing user.authorization.grant for DB call", "user_id", notification.UserID, "error", err)
		return
	} else {
		slog.Debug("Successfully saved user info to DB")
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
		slog.Debug("Working on creating subscription", "user_id", notification.UserID, "type", subscription.Type)
		bodyBytes, err := json.Marshal(subscription)

		if err != nil {
			slog.Warn("Could not marshal subscription for user", "error", err)
		}

		body := string(bodyBytes)
		subType := subscription.Type
		go func() {
			slog.Debug("Inside GO function to create subscription", "user_id", notification.UserID, "type", subType)
			statusCode, response, err := api.CallApi(dbServices, http.MethodPost, "eventsub/subscriptions", body, nil)
			if err != nil {
				slog.Warn("Could not subscribe to event for user", "subscription_type", subType, "error", err, "response", string(response))
			} else {
				slog.Debug("Successfully requested creation of event", "response", string(response), slog.Int("status_code", statusCode))
			}
		}()
	}

	slog.Debug("Publishing event to queue")
	dbServices.Queue.Publish("user.authorization.grant", notification)
}
