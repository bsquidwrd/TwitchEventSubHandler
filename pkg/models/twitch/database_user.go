package twitch

import (
	"database/sql"
	"time"
)

type DatabaseUser struct {
	Id            string              `json:"id"`
	Name          string              `json:"name"`
	Login         string              `json:"login"`
	AvatarUrl     string              `json:"avatar_url,omitempty"`
	Description   string              `json:"description,omitempty"`
	Title         string              `json:"title,omitempty"`
	Language      string              `json:"language,omitempty"`
	CategoryId    string              `json:"category_id,omitempty"`
	CategoryName  string              `json:"category_name,omitempty"`
	LastOnlineAt  sql.Null[time.Time] `json:"last_online_at,omitempty"`
	LastOfflineAt sql.Null[time.Time] `json:"last_offline_at,omitempty"`
	Live          bool                `json:"live,omitempty"`
}
