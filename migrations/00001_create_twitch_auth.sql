-- +goose Up
-- +goose StatementBegin
CREATE TABLE twitch_auth (
	id serial4 NOT NULL,
	client_id varchar NOT NULL,
	access_token varchar NOT NULL,
	expires_at timestamp with time zone NOT NULL,
	expired bool DEFAULT false NOT NULL,
	CONSTRAINT twitch_auth_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE twitch_auth;
-- +goose StatementEnd
