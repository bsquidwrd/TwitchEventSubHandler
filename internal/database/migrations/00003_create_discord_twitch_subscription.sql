-- +goose Up
-- +goose StatementBegin
CREATE TABLE discord_twitch_subscription (
	guild_id varchar NOT NULL,
	user_id varchar NOT NULL,
	webhook_id varchar NOT NULL,
	"token" varchar NOT NULL,
	message varchar DEFAULT '' NOT NULL,
	last_message_id varchar DEFAULT '' NOT NULL,
	last_message_timestamp timestamp with time zone NULL,
	last_online_processed timestamp with time zone NULL,
	CONSTRAINT discord_twitch_subscription_pk PRIMARY KEY (guild_id,user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE discord_twitch_subscription;
-- +goose StatementEnd
