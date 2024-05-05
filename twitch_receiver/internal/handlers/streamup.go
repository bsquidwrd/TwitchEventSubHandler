package handlers

import (
	"context"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/service"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/pkg/models"
)

func processStreamUp(dbServices *service.ReceiverService, notification *models.StreamUpEventSubEvent) {
	slog.Debug("Channel went live", "userid", notification.BroadcasterUserID)

	_, err := dbServices.Database.Exec(context.Background(), `
		update public.twitch_user
		set last_online_at=$2,live=$3
		where id=$1
		`,
		notification.BroadcasterUserID,
		notification.StartedAt,
		true,
	)

	if err != nil {
		slog.Warn("Error processing stream.online for DB call", "user_id", notification.BroadcasterUserID)
		return
	}

	dbServices.Queue.Publish("stream.online", notification)
}
