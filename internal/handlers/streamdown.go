package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processStreamDown(dbServices *database.Service, notification *models.StreamDownEventSubEvent) {
	slog.Info("Channel went offline", "userid", notification.BroadcasterUserID)
	defer dbServices.Queue.Publish("stream.offline", notification)

	go func() {
		_, err := dbServices.Database.Exec(context.Background(), `
		insert into public.twitch_user (id,"name",login,live,last_offline_at)
		values($1,$2,$3,$4,$5)
		on conflict (id) do update
		set "name"=$2,login=$3,live=$4,last_offline_at=$5;
		`,
			notification.BroadcasterUserID,
			notification.BroadcasterUserName,
			notification.BroadcasterUserLogin,
			false,
			time.Now().UTC(),
		)

		if err != nil {
			slog.Warn("Error processing stream.offline for DB call", "userid", notification.BroadcasterUserID)
		}
	}()
}
