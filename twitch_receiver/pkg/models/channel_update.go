// Copied from Twitch CLI
// https://github.com/twitchdev/twitch-cli/blob/83b47aa44a986d3ff47d3800d3fee7983813a7a4/internal/models/channel_update.go
package models

type ChannelUpdateEventSubEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	StreamTitle          string `json:"title"`
	StreamLanguage       string `json:"language"`
	StreamCategoryID     string `json:"category_id"`
	StreamCategoryName   string `json:"category_name"`

	// v1
	IsMature *bool `json:"is_mature,omitempty"`

	// v2
	ContentClassificationLabels []string `json:"content_classification_labels,omitempty"`
}

// Renamed from ChannelUpdateEventSubResponse to ChannelUpdateEventMessage to make more sense in my context
type ChannelUpdateEventMessage struct {
	Subscription EventsubSubscription       `json:"subscription"`
	Event        ChannelUpdateEventSubEvent `json:"event"`
}
