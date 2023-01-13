package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")
)

type ErrInvalidSecretKeySize struct {
	Size int
}

func (e ErrInvalidSecretKeySize) Error() string {
	return fmt.Sprintf("invalid secret key size: %d", e.Size)
}

// Payload is the payload object for the token.
type Payload struct {
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Username  string    `json:"username"`
	ID        uuid.UUID `json:"id"`
	// Expiration is the time at which the token expires
	Expiration int64 `json:"exp"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	payload := &Payload{
		IssuedAt:   now,
		ExpiredAt:  now.Add(duration),
		Username:   username,
		ID:         id,
		Expiration: now.Add(duration).Unix(),
	}
	return payload, nil
}

func (p *Payload) IsExpired() bool {
	return p.ExpiredAt.Before(time.Now())
}

func (p *Payload) Valid() error {
	if p.IsExpired() {
		return ErrExpiredToken
	}
	return nil
}
