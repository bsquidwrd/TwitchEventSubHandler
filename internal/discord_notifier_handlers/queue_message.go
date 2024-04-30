package discordnotifierhandlers

import (
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

	slog.Info("Got message for user", "user_id", userId)
}
