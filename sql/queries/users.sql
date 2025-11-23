-- name: CreateUser :one
insert into users (email)
values ($1)
returning *;
