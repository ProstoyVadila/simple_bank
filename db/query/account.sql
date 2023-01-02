-- name: CreateAccount :one
insert into accounts (
    owner_name,
    balance,
    currency
) values (
    $1, $2, $3
) returning *;

-- name: GetAccount :one
select * from accounts
where id = $1;

-- name: ListAccounts :many
select * from accounts
order by created_at
limit $1
offset $2;

-- name: UpdateAccount :one
update accounts 
set balance = $2
where id = $1
returning *;

-- name: DeleteAccount :exec
delete from accounts
where id = $1;
