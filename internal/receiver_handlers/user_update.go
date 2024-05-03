package receiver_handlers

import (
	"context"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	models "github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
)

func processUserUpdate(dbServices *database.ReceiverService, notification *models.UserUpdateEventSubEvent) {
	slog.Debug("User was updated", "userid", notification.UserID)

	_, err := dbServices.Database.Exec(context.Background(), `
		update public.twitch_user
		set "name"=$2,login=$3,email=$4,email_verified=$5,description=$6
		where id=$1
		`,
		notification.UserID,
		notification.UserName,
		notification.UserLogin,
		notification.Email,
		notification.EmailVerified,
		notification.Description,
	)

	if err != nil {
		slog.Warn("Error processing user.update for DB call", "user_id", notification.UserID)
		return
	}

	dbServices.Queue.Publish("user.update", notification)
}
