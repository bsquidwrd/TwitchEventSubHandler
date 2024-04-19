// Copied from Twitch CLI
// https://github.com/twitchdev/twitch-cli/blob/83b47aa44a986d3ff47d3800d3fee7983813a7a4/internal/models/authorization_revoke.go
package models

// Renamed from AuthorizationRevokeEvent to AuthorizationRevokeEventData so I can have the full event be consistent below
type AuthorizationRevokeEventData struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	ClientID  string `json:"client_id"`
}

// Renamed from AuthorizationRevokeEventSubResponse to AuthorizationRevokeEvent to make more sense in my context
type AuthorizationRevokeEvent struct {
	Subscription EventsubSubscription          `json:"subscription"`
	Event        *AuthorizationRevokeEventData `json:"event,omitempty"`
}
