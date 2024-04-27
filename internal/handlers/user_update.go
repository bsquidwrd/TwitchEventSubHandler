package handlers

import (
	"context"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/models"
)

func processUserUpdate(dbServices *database.Service, notification *models.UserUpdateEventMessage) {
	slog.Info("User was updated", "username", notification.Event.UserName)
	defer dbServices.Queue.Publish("user.update", notification)

	go func() {
		_, err := dbServices.Database.Exec(context.Background(), `
		insert into public.twitch_user (id,"name",login,email,email_verified,description)
		values($1,$2,$3,$4,$5,$6)
		on conflict (id) do update
		set "name"=$2,login=$3,email=$4,email_verified=$5,description=$6;
		`,
			notification.Event.UserID,
			notification.Event.UserName,
			notification.Event.UserLogin,
			notification.Event.Email,
			notification.Event.EmailVerified,
			notification.Event.Description,
		)

		if err != nil {
			slog.Warn("Error processing user.update for DB call", "id", notification.Event.UserID)
		}
	}()
}
