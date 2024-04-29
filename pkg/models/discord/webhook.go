package discord

import (
	"errors"
	"time"
)

const (
	SUPPRESS_EMBEDS        = 1 << 2
	SUPPRESS_NOTIFICATIONS = 1 << 12
)

type WebhookBody struct {
	Id        string  `json:"id,omitempty"`
	Content   string  `json:"content,omitempty"`
	Username  string  `json:"username,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	TTS       bool    `json:"tts,omitempty"`
	Embeds    []embed `json:"embeds,omitempty"`
	Flags     int     `json:"flags,omitempty"`
}

type embed struct {
	Title       string   `json:"title,omitempty"`
	Type        string   `json:"type,omitempty"`
	Description string   `json:"description,omitempty"`
	Url         string   `json:"url,omitempty"`
	Timestamp   string   `json:"timestamp,omitempty"`
	Color       int      `json:"color,omitempty"`
	Footer      footer   `json:"footer,omitempty"`
	Image       image    `json:"image,omitempty"`
	Thumbnail   image    `json:"thumbnail,omitempty"`
	Provider    provider `json:"provider,omitempty"`
	Author      author   `json:"author,omitempty"`
	Fields      []field  `json:"fields,omitempty"`
}

type footer struct {
	Text         string `json:"text"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type image struct {
	Url      string `json:"url"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type provider struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type author struct {
	Name         string `json:"name"`
	Url          string `json:"url,omitempty"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

func (w *WebhookBody) AddEmbed(e embed) error {
	if len(w.Embeds) >= 10 {
		return errors.New("can not have more than 10 embeds")
	}

	w.Embeds = append(w.Embeds, e)
	return nil
}

func NewEmbed(title string, description string, url string, timestamp time.Time, color int) embed {
	return embed{
		Title:       title,
		Type:        "rich",
		Description: description,
		Url:         url,
		Timestamp:   timestamp.Format(time.RFC3339),
		Color:       color,
	}
}

func (e *embed) AddField(name string, value string, inline bool) error {
	if len(e.Fields) >= 25 {
		return errors.New("can not have more than 25 fields in one embed")
	}

	e.Fields = append(e.Fields, field{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
	return nil
}

func (e *embed) WithFooter(text string, iconUrl string) {
	e.Footer = footer{
		Text:    text,
		IconUrl: iconUrl,
	}
}

func (e *embed) WithAuthor(name string, url string, iconUrl string) {
	e.Author = author{
		Name:    name,
		Url:     url,
		IconUrl: iconUrl,
	}
}

func (e *embed) WithProvider(name string, url string) {
	e.Provider = provider{
		Name: name,
		Url:  url,
	}
}

func (e *embed) WithImage(url string) {
	e.Image = image{
		Url: url,
	}
}

func (e *embed) WithThumbnail(url string) {
	e.Thumbnail = image{
		Url: url,
	}
}
