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
		insert into public.twitch_user (id,"name",login,last_online_at,live)
		values($1,$2,$3,$4,$5)
		on conflict (id) do update
		set "name"=$2,login=$3,last_online_at=$4,live=$5;
		`,
		notification.BroadcasterUserID,
		notification.BroadcasterUserName,
		notification.BroadcasterUserLogin,
		notification.StartedAt,
		true,
	)

	if err != nil {
		slog.Warn("Error processing stream.online for DB call", "userid", notification.BroadcasterUserID)
		return
	}

	dbServices.Queue.Publish("stream.online", notification)
}
