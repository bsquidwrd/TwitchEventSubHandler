// Copied from Twitch CLI
// https://github.com/twitchdev/twitch-cli/blob/83b47aa44a986d3ff47d3800d3fee7983813a7a4/internal/models/streamdown.go
package eventsub

// Renamed from StreamDownEventSubResponse to StreamDownEventMessage to make more sense in my context
type StreamDownEventMessage struct {
	Subscription EventsubSubscription    `json:"subscription"`
	Event        StreamDownEventSubEvent `json:"event"`
}

type StreamDownEventSubEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
}
