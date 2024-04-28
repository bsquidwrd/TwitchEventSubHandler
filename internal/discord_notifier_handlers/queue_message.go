package discordnotifierhandlers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ProcessMessage(dbServices *database.DiscordNotifierService, msg amqp.Delivery) {
	var userId string
	switch msg.RoutingKey {
	case "channel.update":
		var event twitch.ChannelUpdateEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		userId = event.BroadcasterUserID
	case "user.update":
		var event twitch.UserUpdateEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		userId = event.UserID
	case "stream.online":
		var event twitch.StreamUpEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		userId = event.BroadcasterUserID
	case "stream.offline":
		var event twitch.StreamDownEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		userId = event.BroadcasterUserID
	}

	if userId == "" {
		return
	}

	dbUser := dbServices.Database.QueryRow(
		context.Background(),
		`
			select id, "name", login, description, title, language, category_id, category_name, last_online_at, last_offline_at, live
			from public.twitch_user
			where id=$1
		`,
		userId,
	)

	var user twitch.DatabaseUser
	err := dbUser.Scan(&user.Id, &user.Name, &user.Login, &user.Description, &user.Title, &user.Language, &user.CategoryId, &user.CategoryName, &user.LastOnlineAt, &user.LastOfflineAt, &user.Live)
	if err != nil {
		slog.Warn("Could not retrieve user from database", err)
		return
	}

	slog.Info("Got user info!", "user", user)
}
