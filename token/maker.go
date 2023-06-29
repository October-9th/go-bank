package token

import "time"

// Maker í an interface for managing tokens
type Maker interface {
	// CreateToken creates and sign new token for specific username and valid duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken check if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
