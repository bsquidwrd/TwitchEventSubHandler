package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_discord_notifier/internal/service"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/pkg/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ProcessMessage(dbServices *service.DiscordNotifierService, msg amqp.Delivery) {
	switch msg.RoutingKey {
	case "channel.update":
		var event models.ChannelUpdateEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		handleChannelUpdate(event)
	case "user.update":
		var event models.UserUpdateEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		handleUserUpdate(event)
	case "stream.online":
		var event models.StreamUpEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		handleStreamOnline(dbServices, event)
	case "stream.offline":
		var event models.StreamDownEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		handleStreamOffline(event)
	}
}
