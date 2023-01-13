package token

import (
	"errors"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto      *paseto.V2
	symetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, ErrInvalidSecretKeySize{Size: chacha20poly1305.KeySize}
	}
	maker := &PasetoMaker{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return p.paseto.Encrypt(p.symetricKey, payload, nil)
}

func (p *PasetoMaker) ValidateToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := p.paseto.Decrypt(token, p.symetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		if errors.Is(err, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	return payload, nil
}
