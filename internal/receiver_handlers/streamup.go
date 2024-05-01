package receiver_handlers

import (
	"context"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	models "github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
)

func processStreamUp(dbServices *database.ReceiverService, notification *models.StreamUpEventSubEvent) {
	slog.Info("Channel went live", "userid", notification.BroadcasterUserID)

	_, err := dbServices.Database.Exec(context.Background(), `
		update public.twitch_user
		set "last_online_at=$2,live=$3
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
