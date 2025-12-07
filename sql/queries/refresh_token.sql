-- name: CreateRefreshToken :one
insert into refresh_tokens(token, user_id, expires_at)
values ($1, $2, $3)
returning *;

-- name: GetRefreshToken :one
select * from refresh_tokens
where token = $1 and revoked_at is null;

-- name: RevokeToken :exec
update refresh_tokens
set revoked_at = current_timestamp, updated_at = current_timestamp
where token = $1 and revoked_at is null;
