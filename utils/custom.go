package utils

import (
	"fmt"

	"github.com/google/uuid"
)

type UUIDString string

type ErrInvalidID struct {
	id string
}

func (e ErrInvalidID) Error() string {
	return fmt.Sprintf("Invalid id format: %v", e.id)
}

// UUID helps to convert ID string to uuid.UUID bc of gin uri params validation bug
func (u *UUIDString) UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(string(*u))
	if id == uuid.Nil {
		return id, ErrInvalidID{id: id.String()}
	}
	return id, err
}
