package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func createRandomHashedPassword(t *testing.T, password string) (string, string) {
	if password == "" {
		password = RandomString(12)
	}
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	require.NotEqual(t, hashedPassword, password)
	return hashedPassword, password
}

func TestPasswordHashFunc(t *testing.T) {
	hashedPass, pass := createRandomHashedPassword(t, "")
	err := CheckPassword(pass, hashedPass)
	require.NoError(t, err)

	wrongPass := RandomString(6)
	err = CheckPassword(wrongPass, hashedPass)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}

func TestHashedPasswordRandomSalt(t *testing.T) {
	pass := RandomString(24)
	hashedPass1, _ := createRandomHashedPassword(t, pass)
	hashedPass2, _ := createRandomHashedPassword(t, pass)
	require.NotEqual(t, hashedPass1, hashedPass2)
}

func TestPasswordHashToCollision(t *testing.T) {
	hashedPass1, _ := createRandomHashedPassword(t, "")
	hashedPass2, _ := createRandomHashedPassword(t, "")
	require.NotEqual(t, hashedPass1, hashedPass2)
}
