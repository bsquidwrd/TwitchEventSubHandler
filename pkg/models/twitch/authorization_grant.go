package twitch

type AuthorizationGrantEvent struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	ClientID  string `json:"client_id"`
}

type AuthorizationGrantEventMessage struct {
	Subscription EventsubSubscription    `json:"subscription"`
	Event        AuthorizationGrantEvent `json:"event,omitempty"`
}
