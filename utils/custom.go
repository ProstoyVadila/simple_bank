package utils

import (
	"reflect"

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

func KindOf(obj interface{}) reflect.Kind {
	val := reflect.ValueOf(obj)
	valType := val.Kind()

	if valType == reflect.Ptr {
		valType = val.Elem().Kind()
	}
	return valType
}
