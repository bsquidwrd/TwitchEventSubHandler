package discordnotifierhandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/discord"
	"github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
)

func escapeMarkdown(input string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`*`, `\*`,
		`_`, `\_`,
		`~`, `\~`,
		"`", "\\`",
		`/`, `\/`,
		`>`, `\>`,
		`|`, `\|`,
	)
	return replacer.Replace(input)
}

func escapeUrl(input string) string {
	return fmt.Sprintf("<%s>", input)
}

func getRelativeTimestamp(timestamp time.Time) string {
	return fmt.Sprintf("<t:%d:R>", timestamp.Unix())
}

func ProcessWebhook(dbServices *database.DiscordNotifierService, userId string) {
	dbUser := dbServices.Database.QueryRow(
		context.Background(),
		`
			select id, "name", login, avatar_url, description, title, language, category_id, category_name, last_online_at, last_offline_at, live
			from public.twitch_user
			where id=$1
		`,
		userId,
	)

	var user twitch.DatabaseUser
	err := dbUser.Scan(&user.Id, &user.Name, &user.Login, &user.AvatarUrl, &user.Description, &user.Title, &user.Language, &user.CategoryId, &user.CategoryName, &user.LastOnlineAt, &user.LastOfflineAt, &user.Live)
	if err != nil {
		slog.Warn("Could not retrieve user from database", err)
		return
	}

	slog.Info("Got user info!", "user", user)

	rows, err := dbServices.Database.Query(
		context.Background(),
		`
			select webhookid, "token", message, last_message_id
			from public.discord_twitch_subscription
			where userid=$1
		`,
		userId,
	)

	if err != nil {
		slog.Warn("Could not get subscriptions from db", err)
	}
	defer rows.Close()

	var subscriptions []discord.Subscription
	for rows.Next() {
		var s discord.Subscription
		err := rows.Scan(&s.WebhookId, &s.Token, &s.Message, &s.LastMessageId)
		if err != nil {
			log.Fatal(err)
		}
		subscriptions = append(subscriptions, s)
	}
	if err := rows.Err(); err != nil {
		slog.Warn("There was an error reading a row", err)
	}

	for _, sub := range subscriptions {
		webhookUrl := sub.GetUrl()
		profileUrl := fmt.Sprintf("https://twitch.tv/%s", user.Login)

		message := sub.Message
		if message == "" {
			message = "{name} is live and playing {game}! {url}"
		}

		title := user.Title
		if title == "" {
			title = "[Not Set]"
		}
		embedColor := 0x9146FF
		if !user.Live {
			embedColor = 0xCFCFCF
		}

		categoryName := user.CategoryName
		if categoryName == "" {
			categoryName = "[Not Set]"
		}

		replacer := strings.NewReplacer(
			"{name}", escapeMarkdown(user.Name),
			"{title}", escapeMarkdown(title),
			"{game}", escapeMarkdown(categoryName),
			"{url}", escapeUrl(profileUrl),
		)
		formattedMessage := replacer.Replace(message)

		if !user.Live {
			offlinePrefix := "**[OFFLINE]**"
			if user.LastOfflineAt.Valid {
				offlineTimestamp := getRelativeTimestamp(user.LastOfflineAt.V)
				offlinePrefix = fmt.Sprintf("%s %s", offlinePrefix, offlineTimestamp)
			}
			formattedMessage = fmt.Sprintf("%s\n%s", offlinePrefix, formattedMessage)
		}

		embed := discord.NewEmbed(
			"",
			escapeMarkdown(user.Title),
			profileUrl,
			user.LastOnlineAt.V,
			embedColor,
		)

		embed.WithAuthor(user.Name, profileUrl, user.AvatarUrl)
		embed.WithThumbnail(user.AvatarUrl)
		embed.AddField("Game", escapeMarkdown(categoryName), true)
		embed.AddField("Stream", profileUrl, true)

		if user.LastOnlineAt.Valid {
			embed.WithFooter("Stream start time", "")
		} else {
			embed.Timestamp = ""
		}

		body := discord.WebhookBody{
			Content: formattedMessage,
		}
		body.AddEmbed(embed)

		go func() {
			data, err := json.Marshal(body)
			if err != nil {
				slog.Warn("Could not marshal body", err)
				return
			}

			requestMethod := http.MethodPost
			if sub.LastMessageId != "" {
				requestMethod = http.MethodPatch
			}

			request, err := http.NewRequest(requestMethod, webhookUrl, bytes.NewReader(data))
			if err != nil {
				slog.Warn("Could not assemble request", err)
				return
			}

			request.Header.Add("Accept", "application/json")
			request.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			var response *http.Response

			// Retry the request 3 times
			for i := 0; i < 3; i++ {
				response, err = client.Do(request)
				if err != nil {
					time.After(1 * time.Second)
					err = nil
					continue
				} else {
					break
				}
			}
			if err != nil {
				slog.Error("Error sending webhook", err)
				return
			}
			defer response.Body.Close()

			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				slog.Error("Unable to parse body from response", err)
				return
			}

			if response.StatusCode == http.StatusNotFound {
				// delete subscription as webhook is not found anymore
				slog.Info("Webhook not found", "body", string(responseBody))
			} else {
				// message sent successfully
				slog.Info("Message sent successfully")
			}
		}()
	}
}
