package e

import (
	"fmt"

	"github.com/google/uuid"
)

type ErrAccountNotFound struct {
	Id uuid.UUID
}

func (e ErrAccountNotFound) Error() string {
	return fmt.Sprintf("Account %v not found", e.Id)
}

type ErrInvalidID struct {
	Id string
}

func (e ErrInvalidID) Error() string {
	return fmt.Sprintf("Invalid id format: %v", e.Id)
}

type ErrInvalidCurrencyType struct {
	ReqCurr     string
	FromAccCurr string
	ToAcctCurr  string
}

func (e ErrInvalidCurrencyType) Error() string {
	return fmt.Sprintf(
		"Mismatched currency types: requested %v, but wanted %v and %v",
		e.ReqCurr,
		e.FromAccCurr,
		e.ToAcctCurr,
	)
}
