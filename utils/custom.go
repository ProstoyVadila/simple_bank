package utils

import "github.com/google/uuid"

type UUIDString string

// UUID helps to convert ID string to uuid.UUID bc of gin validation/convertation bug
func (u *UUIDString) UUID() (uuid.UUID, error) {
	return uuid.Parse(string(*u))
}
