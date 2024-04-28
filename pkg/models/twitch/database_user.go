package twitch

import "time"

type DatabaseUser struct {
	Id            string
	Name          string
	Login         string
	Email         string
	EmailVerified bool
	Description   string
	Title         string
	Language      string
	CategoryID    string
	CategoryName  string
	LastOnlineAt  time.Time
	LastOfflineAt time.Time
	Live          bool
}
