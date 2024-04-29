package discord

import "fmt"

type Subscription struct {
	WebhookId     string `json:"webhookid"`
	Token         string `json:"token"`
	Message       string `json:"message"`
	LastMessageId string `json:last_message_id"`
}

func (s *Subscription) GetUrl() string {
	if s.LastMessageId == "" {
		return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s?wait=true", s.WebhookId, s.Token)
	} else {
		return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s/messages/%s?wait=true", s.WebhookId, s.Token, s.LastMessageId)
	}
}
