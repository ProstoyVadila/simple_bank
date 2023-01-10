-- name: CreateUser :one
insert into users (
    username,
    hashed_password,
    full_name,
    email
) values (
    $1, $2, $3, $4
) returning *;

-- name: GetUser :one
select * from users
where username = $1 limit 1;

