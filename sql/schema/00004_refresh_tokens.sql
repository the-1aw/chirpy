-- +goose Up
-- +goose StatementBegin
create table refresh_tokens(
	token text primary key not null,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	user_id uuid not null references users(id) on delete cascade,
	expires_at timestamp not null,
	revoked_at timestamp default null,
	unique (token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table refresh_tokens;
-- +goose StatementEnd
