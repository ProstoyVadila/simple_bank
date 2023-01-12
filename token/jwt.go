package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	minSecretKeyLength = 32
)

var ErrInvalidSecretKeyLength = fmt.Errorf("it must be at least %v characters long", minSecretKeyLength)

// JWT is a JWT token maker. It implements the Maker interface.
type JWT struct {
	secretKey string
}

// NewJWT creates a new JWT token maker.
func NewJWT(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLength {
		return nil, ErrInvalidSecretKeyLength
	}
	return &JWT{secretKey: secretKey}, nil
}

func (j *JWT) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(j.secretKey))
}

func (j *JWT) ValidateToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
