// Copied from Twitch CLI
// https://github.com/twitchdev/twitch-cli/blob/83b47aa44a986d3ff47d3800d3fee7983813a7a4/internal/models/authorization_revoke.go
package models

type AuthorizationRevokeEvent struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	ClientID  string `json:"client_id"`
}

// Renamed from AuthorizationRevokeEventSubResponse to AuthorizationRevokeEventMessage to make more sense in my context
type AuthorizationRevokeEventMessage struct {
	Subscription EventsubSubscription     `json:"subscription"`
	Event        AuthorizationRevokeEvent `json:"event,omitempty"`
}
