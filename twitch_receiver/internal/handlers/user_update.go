package handlers

import (
	"context"
	"log/slog"

	models "github.com/bsquidwrd/TwitchEventSubHandler/shared/models/eventsub"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_receiver/internal/service"
)

func processUserUpdate(dbServices *service.ReceiverService, notification *models.UserUpdateEventSubEvent) {
	slog.Debug("User was updated", "userid", notification.UserID)

	_, err := dbServices.Database.Exec(context.Background(), `
		update public.twitch_user
		set "name"=$2,login=$3,description=$4
		where id=$1
		`,
		notification.UserID,
		notification.UserName,
		notification.UserLogin,
		notification.Description,
	)

	if err != nil {
		slog.Warn("Error processing user.update for DB call", "user_id", notification.UserID)
		return
	}

	dbServices.Queue.Publish("user.update", notification)
}
