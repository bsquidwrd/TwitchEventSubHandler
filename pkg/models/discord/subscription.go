package discord

import (
	"database/sql"
	"fmt"
	"time"
)

type Subscription struct {
	GuildId              string              `json:"guild_id"`
	UserId               string              `json:"user_id"`
	WebhookId            string              `json:"webhook_id"`
	Token                string              `json:"token"`
	Message              string              `json:"message"`
	LastMessageId        string              `json:"last_message_id"`
	LastMessageTimestamp sql.Null[time.Time] `json:"last_message_timestamp"`
	LastOnlineProcessed  sql.Null[time.Time] `json:"last_online_processed"`
}

func (s *Subscription) GetUrl() string {
	if s.LastMessageId == "" {
		return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s?wait=true", s.WebhookId, s.Token)
	} else {
		return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s/messages/%s?wait=true", s.WebhookId, s.Token, s.LastMessageId)
	}
}
