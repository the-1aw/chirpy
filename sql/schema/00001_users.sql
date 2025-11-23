-- +goose Up
-- +goose StatementBegin
create table users(
	id uuid primary key default gen_random_uuid(),
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	email text not null,
	unique(email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
