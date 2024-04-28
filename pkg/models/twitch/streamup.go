// Copied from Twitch CLI
// https://github.com/twitchdev/twitch-cli/blob/83b47aa44a986d3ff47d3800d3fee7983813a7a4/internal/models/streamup.go
package models

// Renamed from StreamUpEventSubResponse to StreamUpEventMessage to make more sense in my context
type StreamUpEventMessage struct {
	Subscription EventsubSubscription  `json:"subscription"`
	Event        StreamUpEventSubEvent `json:"event"`
}

type StreamUpEventSubEvent struct {
	ID                   string `json:"id"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Type                 string `json:"type"`
	StartedAt            string `json:"started_at"`
}
