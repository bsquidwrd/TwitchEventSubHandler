package discordnotifierhandlers

import (
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
)

func handleChannelUpdate(event twitch.ChannelUpdateEventSubEvent) {
	slog.Debug("Handling channel.update", "user_id", event.BroadcasterUserID)
}

func handleUserUpdate(event twitch.UserUpdateEventSubEvent) {
	slog.Debug("Handling user.update", "user_id", event.UserID)
}

func handleStreamOnline(event twitch.StreamUpEventSubEvent) {
	slog.Debug("Handling stream.online", "user_id", event.BroadcasterUserID)
}

func handleStreamOffline(event twitch.StreamDownEventSubEvent) {
	slog.Debug("Handling stream.offline", "user_id", event.BroadcasterUserID)
}
