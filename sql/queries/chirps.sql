-- name: CreateChirp :one
insert into chirps (user_id, body)
values ($1, $2)
returning *;

-- name: GetChirps :many
select * from chirps
order by created_at;

-- name: GetChirpById :one
select * from chirps
where id = $1;

-- name: DeleteAllChirps :exec
delete from chirps;
