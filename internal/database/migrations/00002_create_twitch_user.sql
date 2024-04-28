-- +goose Up
-- +goose StatementBegin
CREATE TABLE twitch_user (
	id varchar NOT NULL,
	"name" varchar NOT NULL,
	login varchar NOT NULL,
	email varchar DEFAULT '' NOT NULL,
	email_verified boolean DEFAULT false NOT NULL,
	description varchar DEFAULT '' NOT NULL,
	title varchar DEFAULT '' NOT NULL,
	"language" varchar DEFAULT '' NOT NULL,
	category_id varchar DEFAULT '' NOT NULL,
	category_name varchar DEFAULT '' NOT NULL,
	last_online_at timestamp with time zone NULL,
	last_offline_at timestamp with time zone NULL,
	live boolean DEFAULT false NOT NULL,
	CONSTRAINT twitch_user_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE twitch_user;
-- +goose StatementEnd
