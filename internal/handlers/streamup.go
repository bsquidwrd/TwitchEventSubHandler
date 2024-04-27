package handlers

import (
	"context"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processStreamUp(dbServices *database.Service, notification *models.StreamUpEventMessage) {
	slog.Info("Channel went live", "username", notification.Event.BroadcasterUserName)

	if notification.Event.Type != "live" {
		return
	}
	defer dbServices.Queue.Publish("stream.online", notification)

	go dbServices.Database.Exec(context.Background(), `
		insert into public.twitch_user (id,"name",login,last_online_at,live)
		values($1,$2,$3,$4,$5)
		on conflict (id) do update
		set "name"=$2,login=$3,last_online_at=$4,live=$5;
		`,
		notification.Event.BroadcasterUserID,
		notification.Event.BroadcasterUserName,
		notification.Event.BroadcasterUserLogin,
		notification.Event.StartedAt,
		true,
	)
}
