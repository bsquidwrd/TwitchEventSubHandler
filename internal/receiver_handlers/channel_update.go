package receiver_handlers

import (
	"context"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	models "github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
)

func processChannelUpdate(dbServices *database.ReceiverService, notification *models.ChannelUpdateEventSubEvent) {
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
