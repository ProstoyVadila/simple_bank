-- name: CreateTransfer :one
insert into transfers (
    from_account_id,
    to_account_id,
    amount
) values (
    $1, $2, $3
) returning *;

-- name: GetTrasfer :one
select * from transfers
where id = $1;

-- name: GetTransfersByFromAccount :many
select * from transfers
where from_account_id = $1
order by created_at;

-- name: GetTransfersByToAccount :many
select * from transfers
where to_account_id = $1
order by created_at;


-- name: ListTransfers :many
select * from transfers
order by created_at
limit $1
offset $2;

