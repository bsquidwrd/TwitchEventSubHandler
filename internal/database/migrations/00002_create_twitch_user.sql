-- +goose Up
-- +goose StatementBegin
CREATE TABLE twitch_user (
	id varchar NOT NULL,
	"name" varchar NULL,
	login varchar NULL,
	email varchar NULL,
	email_verified boolean DEFAULT false NOT NULL,
	description varchar NULL,
	title varchar NULL,
	"language" varchar NULL,
	category_id varchar NULL,
	category_name varchar NULL,
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
