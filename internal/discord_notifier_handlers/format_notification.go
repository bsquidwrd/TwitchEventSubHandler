package discordnotifierhandlers

import (
	"fmt"
	"strings"
	"time"

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

func getNotificationPayload(user twitch.DatabaseUser, subscription discord.Subscription) discord.WebhookBody {
	profileUrl := fmt.Sprintf("https://twitch.tv/%s", user.Login)

	message := subscription.Message
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

	return body
}
