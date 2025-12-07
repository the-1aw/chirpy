-- +goose Up
-- +goose StatementBegin
alter table users add is_chirpy_red boolean not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users drop column is_chirpy_red;
-- +goose StatementEnd
