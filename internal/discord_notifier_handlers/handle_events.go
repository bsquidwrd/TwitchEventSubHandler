package discordnotifierhandlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
)

func handleChannelUpdate(event twitch.ChannelUpdateEventSubEvent) {
	slog.Debug("Handling channel.update", "user_id", event.BroadcasterUserID)
}

func handleUserUpdate(event twitch.UserUpdateEventSubEvent) {
	slog.Debug("Handling user.update", "user_id", event.UserID)
}

func handleStreamOffline(event twitch.StreamDownEventSubEvent) {
	slog.Debug("Handling stream.offline", "user_id", event.BroadcasterUserID)
}

func handleStreamOnline(dbServices *database.DiscordNotifierService, event twitch.StreamUpEventSubEvent) {
	slog.Debug("Handling stream.online", "user_id", event.BroadcasterUserID)

	dbUser := dbServices.Database.QueryRow(
		context.Background(),
		`
			select
				id,"name",login,avatar_url,email,email_verified,description,title,"language"
				,category_id,category_name,last_online_at,last_offline_at,live
			from public.twitch_user
			where id=$1
		`,
		event.BroadcasterUserID,
	)

	var user twitch.DatabaseUser
	dbUser.Scan(
		&user.Id, &user.Name, &user.Login, &user.AvatarUrl, &user.Email, &user.EmailVerified,
		&user.Description, &user.Title, &user.Language, &user.CategoryId, &user.CategoryName,
		&user.LastOnlineAt, &user.LastOfflineAt, &user.Live,
	)

	subscriptions := getUserSubscriptions(dbServices, user.Id)

	for _, sub := range subscriptions {
		inCooldown := false

		if sub.LastOnlineProcessed.Valid {
			if time.Since(sub.LastOnlineProcessed.V) <= 1*time.Hour {
				inCooldown = true
			}
		}

		if !inCooldown {
			slog.Debug("Not in cooldown for subscription", "guild_id", sub.GuildId, "user_id", sub.UserId)
		}

		streamStartedAt, _ := time.Parse(time.RFC3339, event.StartedAt)

		dbServices.Database.Exec(
			context.Background(),
			`
				update public.discord_twitch_subscription
				set last_message_timestamp=$3,last_online_processed=$4
				where guild_id=$1 and user_id=$2
			`,
			sub.GuildId,
			sub.UserId,
			time.Now(),
			streamStartedAt,
		)
	}
}
