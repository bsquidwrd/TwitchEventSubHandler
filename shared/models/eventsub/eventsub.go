// Copied from Twitch CLI
// https://github.com/twitchdev/twitch-cli/blob/83b47aa44a986d3ff47d3800d3fee7983813a7a4/internal/models/eventsub.go
package eventsub

type EventsubSubscription struct {
	ID        string            `json:"id,omitempty"`
	Status    string            `json:"status,omitempty"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition EventsubCondition `json:"condition"`
	Transport EventsubTransport `json:"transport"`
	CreatedAt string            `json:"created_at,omitempty"`
	Cost      int64             `json:"cost,omitempty"`
}

type EventsubTransport struct {
	Method    string `json:"method"`
	Callback  string `json:"callback,omitempty"`
	SessionID string `json:"session_id,omitempty"`
	Secret    string `json:"secret,omitempty"`
}

type EventsubCondition struct {
	BroadcasterUserID     string `json:"broadcaster_user_id,omitempty"`
	ToBroadcasterUserID   string `json:"to_broadcaster_user_id,omitempty"`
	UserID                string `json:"user_id,omitempty"`
	FromBroadcasterUserID string `json:"from_broadcaster_user_id,omitempty"`
	ModeratorUserID       string `json:"moderator_user_id,omitempty"`
	ClientID              string `json:"client_id,omitempty"`
	ExtensionClientID     string `json:"extension_client_id,omitempty"`
	OrganizationID        string `json:"organization_id,omitempty"`
	CategoryID            string `json:"category_id,omitempty"`
	CampaignID            string `json:"campaign_id,omitempty"`
}

// Renamed from EventsubResponse to EventsubMessage to make more sense in this context
type EventsubMessage struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        interface{}          `json:"event,omitempty"`
}

type EventsubSubscriptionVerification struct {
	Challenge    string               `json:"challenge"`
	Subscription EventsubSubscription `json:"subscription"`
}

type EventsubSubscriptionList struct {
	Total        int                    `json:"total"`
	Data         []EventsubSubscription `json:"data"`
	TotalCost    int                    `json:"total_cost"`
	MaxTotalCost int                    `json:"max_total_cost"`
}
