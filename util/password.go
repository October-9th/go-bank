package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword will compute the bcrypte hash password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hashed password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword will check the provided password and the hashed password is correct or not
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
