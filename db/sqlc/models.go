// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	CreatedAt time.Time `json:"created_at"`
	OwnerName string    `json:"owner_name"`
	Currency  string    `json:"currency"`
	Balance   int64     `json:"balance"`
	ID        uuid.UUID `json:"id"`
}

type Entry struct {
	CreatedAt time.Time `json:"created_at"`
	Amount    int64     `json:"amount"`
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
}

type Transfer struct {
	CreatedAt     time.Time `json:"created_at"`
	Amount        int64     `json:"amount"`
	ID            uuid.UUID `json:"id"`
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
}

type User struct {
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
}
