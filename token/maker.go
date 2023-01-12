package token

import "time"

// Maker is an interface for managing tokens.
type Maker interface {
	// CreateToken creates a new token for the given username and valid duration.
	CreateToken(username string, duration time.Duration) (string, error)
	// ValidateToken validates the given token and returns the payload object.
	ValidateToken(token string) (*Payload, error)
}
