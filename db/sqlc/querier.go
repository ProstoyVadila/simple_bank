// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
	DeleteEntriesByAccount(ctx context.Context, accountID uuid.UUID) error
	DeleteEntry(ctx context.Context, id uuid.UUID) error
	GetAccount(ctx context.Context, id uuid.UUID) (Account, error)
	GetAccountForUpdate(ctx context.Context, id uuid.UUID) (Account, error)
	GetEntriesByAccount(ctx context.Context, accountID uuid.UUID) ([]Entry, error)
	GetEntry(ctx context.Context, id uuid.UUID) (Entry, error)
	GetTransfersByFromAccount(ctx context.Context, fromAccountID uuid.UUID) ([]Transfer, error)
	GetTransfersByToAccount(ctx context.Context, toAccountID uuid.UUID) ([]Transfer, error)
	GetTrasfer(ctx context.Context, id uuid.UUID) (Transfer, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	UpdateEntry(ctx context.Context, arg UpdateEntryParams) (Entry, error)
}

var _ Querier = (*Queries)(nil)
