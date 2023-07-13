package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

// JWTMaker is a JSON web token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker create a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: %d, must be at least %d characters", len(secretKey), minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayLoad(username, duration)
	if err != nil {
		return "", payload, err
	}

	// Create token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	// KeyFunc recevied parsed but unsigh token and then verify its header
	// to check if its matches the signing algorithm used to sign the token
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Get the signing algorithm
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		// The signing algorithm doesn't match the the token algorithm
		if !ok {
			return nil, ErrInvalidToken
		}

		// Else return the secret key that used to sign the token
		return []byte(maker.secretKey), nil
	}

	// Parse token
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		// Figure the error type converting the error to the jwt.ValidationError
		verr, ok := err.(*jwt.ValidationError)
		// if the convertion is ok
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// converting jwt.Token.Claims into payload obj
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil

}
