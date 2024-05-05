package handlers

import (
	"context"
	"log/slog"

	models "github.com/bsquidwrd/TwitchEventSubHandler/shared/models/eventsub"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/service"
)

func processChannelUpdate(dbServices *service.ReceiverService, notification *models.ChannelUpdateEventSubEvent) {
	slog.Debug("Channel was updated", "userid", notification.BroadcasterUserID)

	_, err := dbServices.Database.Exec(context.Background(), `
		update public.twitch_user
		set title=$2,"language"=$3,category_id=$4,category_name=$5
		where id=$1
		`,
		notification.BroadcasterUserID,
		notification.StreamTitle,
		notification.StreamLanguage,
		notification.StreamCategoryID,
		notification.StreamCategoryName,
	)

	if err != nil {
		slog.Warn("Error processing channel.update for DB call", "user_id", notification.BroadcasterUserID)
		return
	}

	dbServices.Queue.Publish("channel.update", notification)
}
