package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	dbModels "github.com/bsquidwrd/TwitchEventSubHandler/shared/models/database"
	models "github.com/bsquidwrd/TwitchEventSubHandler/shared/models/eventsub"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_discord_notifier/internal/service"
)

func handleChannelUpdate(event models.ChannelUpdateEventSubEvent) {
	slog.Debug("Handling channel.update", "user_id", event.BroadcasterUserID)
}

func handleUserUpdate(event models.UserUpdateEventSubEvent) {
	slog.Debug("Handling user.update", "user_id", event.UserID)
}

func handleStreamOffline(event models.StreamDownEventSubEvent) {
	slog.Debug("Handling stream.offline", "user_id", event.BroadcasterUserID)
}

func handleStreamOnline(dbServices *service.DiscordNotifierService, event models.StreamUpEventSubEvent) {
	slog.Debug("Handling stream.online", "user_id", event.BroadcasterUserID)

	dbUser := dbServices.Database.QueryRow(
		context.Background(),
		`
			select
				id,"name",login,avatar_url,description,title,"language"
				,category_id,category_name,last_online_at,last_offline_at,live
			from public.twitch_user
			where id=$1
		`,
		event.BroadcasterUserID,
	)

	var user dbModels.TwitchUser
	dbUser.Scan(
		&user.Id, &user.Name, &user.Login, &user.AvatarUrl, &user.Description, &user.Title,
		&user.Language, &user.CategoryId, &user.CategoryName,
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
			webhookUrl := sub.GetUrl()
			payload := getNotificationPayload(user, sub)

			statusCode, response, err := executeWebhook(webhookUrl, http.MethodPost, payload)
			if err != nil {
				slog.Warn("Error executing webhook for subscription", "guild_id", sub.GuildId, "user_id", sub.UserId, "error", err)
			}

			if statusCode == http.StatusNotFound {
				deleteSubscription(dbServices, sub.GuildId, sub.UserId)
				continue
			}

			dbServices.Database.Exec(
				context.Background(),
				`
					update public.discord_twitch_subscription
					set last_message_id=$3,last_message_timestamp=$4
					where guild_id=$1 and user_id=$2
				`,
				sub.GuildId,
				sub.UserId,
				response.Id,
				time.Now(),
			)
		}

		streamStartedAt, _ := time.Parse(time.RFC3339, event.StartedAt)

		dbServices.Database.Exec(
			context.Background(),
			`
				update public.discord_twitch_subscription
				set last_online_processed=$3
				where guild_id=$1 and user_id=$2
			`,
			sub.GuildId,
			sub.UserId,
			streamStartedAt,
		)
	}
}
