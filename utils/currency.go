package utils

import (
	"fmt"

	"github.com/ProstoyVadila/simple_bank/e"
)

const (
	USD = "USD"
	KZT = "KZT"
	PHP = "PHP"
	EUR = "EUR"
)

var Currencies = [4]string{USD, KZT, PHP, EUR}

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, KZT, PHP, EUR:
		return true
	default:
		return false
	}
}

// IsValidCurrency checks a same currency type for all participants in a money transfer
func ValidateCurrency(requestedCurrency, fromAccountCurrency, toAccountCurrency string) error {
	if !IsSupportedCurrency(requestedCurrency) {
		return e.ErrInvalidCurrencyType{
			Curr: requestedCurrency,
			Msg:  fmt.Sprintf("Unsupported currency type: %v, supported: %v", requestedCurrency, Currencies),
		}
	}
	if !(requestedCurrency == fromAccountCurrency && requestedCurrency == toAccountCurrency) {
		return e.ErrInvalidCurrencyType{
			Curr: requestedCurrency,
			Msg: fmt.Sprintf(
				"Mismatched currency types: requested %v, but wanted %v and %v",
				requestedCurrency,
				fromAccountCurrency,
				toAccountCurrency,
			),
		}
	}
	return nil
}
