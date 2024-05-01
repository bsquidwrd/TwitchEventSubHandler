package receiver_handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	models "github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
)

func processStreamDown(dbServices *database.ReceiverService, notification *models.StreamDownEventSubEvent) {
	slog.Info("Channel went offline", "userid", notification.BroadcasterUserID)

	_, err := dbServices.Database.Exec(context.Background(), `
		insert into public.twitch_user (id,live,last_offline_at)
		values($1,$2,$3)
		on conflict (id) do update
		set live=$2,last_offline_at=$3;
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
