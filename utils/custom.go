package utils

import (
	"github.com/ProstoyVadila/simple_bank/e"
	"github.com/google/uuid"
)

type UUIDString string

// UUID helps to convert ID string to uuid.UUID bc of gin uri params validation bug
func (u *UUIDString) UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(string(*u))
	if id == uuid.Nil {
		return id, e.ErrInvalidID{Id: id.String()}
	}
	return id, err
}

// IsValidCurrency checks a same currency type for all participants in a money transfer
func ValidateCurrency(requestedCurrency, fromAccountCurrency, toAccountCurrency string) error {
	if !(requestedCurrency == fromAccountCurrency && requestedCurrency == toAccountCurrency) {
		return e.ErrInvalidCurrencyType{
			ReqCurr:     requestedCurrency,
			FromAccCurr: fromAccountCurrency,
			ToAcctCurr:  toAccountCurrency,
		}
	}
	return nil
}
