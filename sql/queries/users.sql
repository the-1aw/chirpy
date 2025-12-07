-- name: CreateUser :one
insert into users (email, hashed_password)
values ($1, $2)
returning *;

-- name: GetUserByEmail :one
select * from users
where email = $1;

-- name: UpdateUserById :exec
update users
set hashed_password = $2, email = $3, updated_at = current_timestamp
where id = $1
returning *;

-- name: DeleteAllUsers :exec
delete from users;
