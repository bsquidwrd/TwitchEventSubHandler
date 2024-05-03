package receiver_handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	models "github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
)

func processStreamDown(dbServices *database.ReceiverService, notification *models.StreamDownEventSubEvent) {
	slog.Debug("Channel went offline", "userid", notification.BroadcasterUserID)

	_, err := dbServices.Database.Exec(context.Background(), `
		update public.twitch_user
		set live=$2,last_offline_at=$3
		where id=$1
		`,
		notification.BroadcasterUserID,
		false,
		time.Now().UTC(),
	)

	if err != nil {
		slog.Warn("Error processing stream.offline for DB call", "user_id", notification.BroadcasterUserID)
		return
	}

	dbServices.Queue.Publish("stream.offline", notification)
}
