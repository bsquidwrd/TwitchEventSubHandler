package handlers

import (
	"context"
	"log"
	"log/slog"

	twitch "github.com/bsquidwrd/TwitchEventSubHandler/shared/models/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_discord_notifier/internal/service"
	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_discord_notifier/pkg/models"
)

func getUserSubscriptions(dbServices *service.DiscordNotifierService, userId string) []models.Subscription {
	dbUser := dbServices.Database.QueryRow(
		context.Background(),
		`
			select id, "name", login, avatar_url, description, title, language, category_id, category_name, last_online_at, last_offline_at, live
			from public.twitch_user
			where id=$1
		`,
		userId,
	)

	var user twitch.TwitchUser
	err := dbUser.Scan(&user.Id, &user.Name, &user.Login, &user.AvatarUrl, &user.Description, &user.Title, &user.Language, &user.CategoryId, &user.CategoryName, &user.LastOnlineAt, &user.LastOfflineAt, &user.Live)
	if err != nil {
		slog.Warn("Could not retrieve user from database", err)
		return nil
	}

	rows, err := dbServices.Database.Query(
		context.Background(),
		`
			select guild_id, user_id, webhook_id, "token", message, last_message_id, last_message_timestamp, last_online_processed
			from public.discord_twitch_subscription
			where user_id=$1
		`,
		userId,
	)

	if err != nil {
		slog.Warn("Could not get subscriptions from db", err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var s models.Subscription
		err := rows.Scan(&s.GuildId, &s.UserId, &s.WebhookId, &s.Token, &s.Message, &s.LastMessageId, &s.LastMessageTimestamp, &s.LastOnlineProcessed)
		if err != nil {
			log.Fatal(err)
		}
		subscriptions = append(subscriptions, s)
	}
	if err := rows.Err(); err != nil {
		slog.Warn("There was an error reading a row", err)
	}

	return subscriptions
}

func deleteSubscription(dbServices *service.DiscordNotifierService, guildId string, userId string) {
	_, _ = dbServices.Database.Exec(
		context.Background(),
		`
			delete from public.discord_twitch_subscription
			where guild_id=$1 and user_id=$1
		`,
		guildId,
		userId,
	)
}
