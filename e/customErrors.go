package e

import (
	"fmt"
	"net/http"
)

type ErrUnauthorized struct {
	Msg string
	Obj string
}

func (e ErrUnauthorized) Error() string {
	if e.Obj == "" && e.Msg == "" {
		return "Unauthorized"
	} else if e.Obj == "" {
		return e.Msg
	}
	return fmt.Sprintf(e.Msg, e.Obj)
}

func (e ErrUnauthorized) StatusCode() int {
	return http.StatusUnauthorized
}

type ErrEntityNotFound struct {
	EntityName string
}

func (e ErrEntityNotFound) Error() string {
	return fmt.Sprintf("%v not found", e.EntityName)
}

type ErrInvalidUUID struct {
	Id string
}

func (e ErrInvalidUUID) Error() string {
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
