package handlers

import (
	"context"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processChannelUpdate(dbServices *database.Service, notification *models.ChannelUpdateEventSubEvent) {
	slog.Info("Channel was updated", "userid", notification.BroadcasterUserID)
	defer dbServices.Queue.Publish("channel.update", notification)

	go func() {
		_, err := dbServices.Database.Exec(context.Background(), `
		insert into public.twitch_user (id,"name",login,title,"language",category_id,category_name)
		values($1,$2,$3,$4,$5,$6,$7)
		on conflict (id) do update
		set "name"=$2,login=$3,title=$4,"language"=$5,category_id=$6,category_name=$7;
		`,
			notification.BroadcasterUserID,
			notification.BroadcasterUserName,
			notification.BroadcasterUserLogin,
			notification.StreamTitle,
			notification.StreamLanguage,
			notification.StreamCategoryID,
			notification.StreamCategoryName,
		)

		if err != nil {
			slog.Warn("Error processing channel.update for DB call", "userid", notification.BroadcasterUserID)
		}
	}()
}
