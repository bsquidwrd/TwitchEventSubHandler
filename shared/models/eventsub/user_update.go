// Copied from Twitch CLI
// https://github.com/twitchdev/twitch-cli/blob/83b47aa44a986d3ff47d3800d3fee7983813a7a4/internal/models/eventsub.go
package eventsub

// Renamed from UserUpdateEventSubResponse to UserUpdateEventMessage to make more sense in my context
type UserUpdateEventMessage struct {
	Subscription EventsubSubscription    `json:"subscription"`
	Event        UserUpdateEventSubEvent `json:"event"`
}

type UserUpdateEventSubEvent struct {
	UserID        string `json:"user_id"`
	UserLogin     string `json:"user_login"`
	UserName      string `json:"user_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Description   string `json:"description"`
}
