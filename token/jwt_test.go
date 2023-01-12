package token

import (
	"testing"
	"time"

	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

type testUserToken struct {
	maker    Maker
	duration time.Duration
	username string
	token    string
}

func createUserToken(t *testing.T, username string, duration time.Duration) testUserToken {
	maker, err := NewJWT(utils.RandomString(32))
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

func TestJWTMaker(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		testFunc  func(t *testing.T, maker Maker)
	}{
		{
			name:      "valid",
			secretKey: utils.RandomString(32),
			testFunc: func(t *testing.T, maker Maker) {
				require.NotNil(t, maker)
			},
		},
		{
			name:      "invalid",
			secretKey: utils.RandomString(16),
			testFunc: func(t *testing.T, maker Maker) {
				require.Nil(t, maker)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maker, err := NewJWT(tt.secretKey)
			if err != nil {
				require.Equal(t, ErrInvalidSecretKeyLength, err)
			}
			tt.testFunc(t, maker)
		})
	}

}

func TestCreateJWTToken(t *testing.T) {
	createUserToken(t, utils.RandomOwner(), time.Minute)
}

func TestValidateJWTToken(t *testing.T) {
	tests := []struct {
		name     string
		user     testUserToken
		testFunc func(t *testing.T, userToken testUserToken)
	}{
		{
			name: "valid token",
			user: createUserToken(t, utils.RandomOwner(), time.Minute),
			testFunc: func(t *testing.T, userToken testUserToken) {
				payload, err := userToken.maker.ValidateToken(userToken.token)
				require.NoError(t, err)
				require.NotNil(t, payload)
				require.Equal(t, userToken.username, payload.Username)
			},
		},
		{
			name: "expired token",
			user: createUserToken(t, utils.RandomOwner(), time.Nanosecond),
			testFunc: func(t *testing.T, userToken testUserToken) {
				time.Sleep(time.Nanosecond * 5)
				payload, err := userToken.maker.ValidateToken(userToken.token)
				require.Error(t, err)
				require.Nil(t, payload)
				require.EqualError(t, err, ErrExpiredToken.Error())
			},
		},
		{
			name: "invalid token length",
			user: createUserToken(t, utils.RandomOwner(), time.Minute),
			testFunc: func(t *testing.T, userToken testUserToken) {
				payload, err := userToken.maker.ValidateToken("invalid token")
				require.Error(t, err)
				require.Nil(t, payload)
				require.EqualError(t, err, ErrInvalidToken.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.user)
		})
	}
}
