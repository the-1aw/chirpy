-- +goose Up
-- +goose StatementBegin
create table chirps(
	id uuid primary key default gen_random_uuid(),
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	body text not null,
	user_id uuid not null references users(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chirps;
-- +goose StatementEnd
