package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Store provides an execution of queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides an execution of queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db, Queries: New(db)}
}

// execTx executes a func within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// TODO(vadim): figure out with isolation levels in sqlc framework
	tx, err := store.db.BeginTx(
		ctx,
		// &sql.TxOptions{Isolation: ...},
		nil,
	)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback error; %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input params if the transfer tansaction
type TransferTxParams struct {
	Currency    string  `json:"currency"`
	Amount      int64   `json:"amount"`
	FromAccount Account `json:"from_account"`
	ToAccount   Account `json:"to_account"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
	Transfer    Transfer `json:"transfer"`
	Currency    string   `json:"currency"`
}

// TransferTx provides a money transfer between accounts.
// It creates a transfer record, adds account etries, and updates accounts' balance in one transaction.
func (store *SQLStore) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Currency = args.Currency

		// Creating a entries records
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccount.ID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccount.ID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		//  TODO: find out a better way
		// Avoiding the db deadlock on a cuncurrent "A -> B, B -> A" selecting/updateting
		if args.FromAccount.CreatedAt.Before(args.ToAccount.CreatedAt) {
			result.FromAccount, result.ToAccount, err = transferMoney(
				ctx,
				q,
				args.FromAccount.ID,
				args.ToAccount.ID,
				-args.Amount,
				args.Amount,
			)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = transferMoney(
				ctx,
				q,
				args.ToAccount.ID,
				args.FromAccount.ID,
				args.Amount,
				-args.Amount,
			)
			if err != nil {
				return err
			}
		}

		// Creating a trasfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccount.ID,
			ToAccountID:   args.ToAccount.ID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	account1ID,
	account2ID uuid.UUID,
	amount1,
	amount2 int64,
) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account1ID,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2ID,
		Amount: amount2,
	})
	return
}
