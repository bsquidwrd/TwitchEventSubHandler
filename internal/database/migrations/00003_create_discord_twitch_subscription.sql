-- +goose Up
-- +goose StatementBegin
CREATE TABLE discord_twitch_subscription (
	guildid varchar NOT NULL,
	userid varchar NOT NULL,
	webhookid varchar NOT NULL,
	"token" varchar NOT NULL,
	message varchar DEFAULT '' NOT NULL,
	last_message_id varchar DEFAULT '' NOT NULL,
	CONSTRAINT discord_twitch_subscription_pk PRIMARY KEY (guildid,userid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP discord_twitch_subscription twitch_user;
-- +goose StatementEnd
