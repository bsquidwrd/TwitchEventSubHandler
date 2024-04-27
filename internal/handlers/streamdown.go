package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processStreamDown(dbServices *database.Service, notification *models.StreamDownEventMessage) {
	slog.Info("Channel went offline", "username", notification.Event.BroadcasterUserName)
	defer dbServices.Queue.Publish("stream.offline", notification)

	go dbServices.Database.Exec(context.Background(), `
		insert into public.twitch_user (id,"name",login,live,last_offline_at)
		values($1,$2,$3,$4,$5)
		on conflict (id) do update
		set "name"=$2,login=$3,live=$4,last_offline_at=$5;
		`,
		notification.Event.BroadcasterUserID,
		notification.Event.BroadcasterUserName,
		notification.Event.BroadcasterUserLogin,
		false,
		time.Now().UTC(),
	)
}
