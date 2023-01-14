package e

import (
	"fmt"
)

type ErrEntityNotFound struct {
	EntityName string
}

func (e ErrEntityNotFound) Error() string {
	return fmt.Sprintf("%v not found", e.EntityName)
}

type ErrInvalidID struct {
	Id string
}

func (e ErrInvalidID) Error() string {
	return fmt.Sprintf("Invalid id format: %v", e.Id)
}

type ErrInvalidCurrencyType struct {
	Curr string
	Msg  string
}

func (e ErrInvalidCurrencyType) Error() string {
	if e.Msg == "" {
		return fmt.Sprintf("Invalid currency type: %v", e.Curr)
	}
	return e.Msg
}

type ErrThrottling struct {
	Msg string
}

func (e ErrThrottling) Error() string {
	if e.Msg == "" {
		return fmt.Sprintln("Limit exceeded")
	}
	return e.Msg
}
