package token

import (
	"testing"
	"time"

	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func createUserPasetoToken(t *testing.T, username string, duration time.Duration) testUserToken {
	maker, err := NewPaseto(utils.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return testUserToken{
		maker:    maker,
		duration: duration,
		username: username,
		token:    token,
	}
}

func TestCreatePasetoToken(t *testing.T) {
	createUserPasetoToken(t, utils.RandomOwner(), time.Minute)
}

func TestValidatePasetoToken(t *testing.T) {
	tests := []struct {
		testFunc func(t *testing.T, userToken testUserToken)
		user     testUserToken
		name     string
	}{
		{
			name: "valid token",
			user: createUserPasetoToken(t, utils.RandomOwner(), time.Minute),
			testFunc: func(t *testing.T, userToken testUserToken) {
				payload, err := userToken.maker.ValidateToken(userToken.token)
				require.NoError(t, err)
				require.NotNil(t, payload)
				require.Equal(t, userToken.username, payload.Username)
			},
		},
		{
			name: "expired token",
			user: createUserPasetoToken(t, utils.RandomOwner(), time.Nanosecond),
			testFunc: func(t *testing.T, userToken testUserToken) {
				time.Sleep(time.Nanosecond * 5)
				payload, err := userToken.maker.ValidateToken(userToken.token)
				require.Error(t, err)
				require.Nil(t, payload)
			},
		},
		{
			name: "invalid token",
			user: createUserPasetoToken(t, utils.RandomOwner(), time.Minute),
			testFunc: func(t *testing.T, userToken testUserToken) {
				payload, err := userToken.maker.ValidateToken(userToken.token + "invalid")
				require.Error(t, err)
				require.Nil(t, payload)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.user)
		})
	}
}
