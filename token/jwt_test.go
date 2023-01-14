package token

import (
	"testing"
	"time"

	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

type testUserToken struct {
	maker    Maker
	username string
	token    string
	duration time.Duration
}

func createUserJWTToken(t *testing.T, username string, duration time.Duration) testUserToken {
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

func TestCreateJWTToken(t *testing.T) {
	createUserJWTToken(t, utils.RandomOwner(), time.Minute)
}

func TestValidateJWTToken(t *testing.T) {
	tests := []struct {
		testFunc func(t *testing.T, userToken testUserToken)
		user     testUserToken
		name     string
	}{
		{
			name: "valid token",
			user: createUserJWTToken(t, utils.RandomOwner(), time.Minute),
			testFunc: func(t *testing.T, userToken testUserToken) {
				payload, err := userToken.maker.ValidateToken(userToken.token)
				require.NoError(t, err)
				require.NotNil(t, payload)
				require.Equal(t, userToken.username, payload.Username)
			},
		},
		{
			name: "expired token",
			user: createUserJWTToken(t, utils.RandomOwner(), time.Nanosecond),
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
			user: createUserJWTToken(t, utils.RandomOwner(), time.Minute),
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

func TestInvalidTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWT(utils.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.ValidateToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
