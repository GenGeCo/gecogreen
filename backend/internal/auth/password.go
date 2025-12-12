package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// bcrypt cost - higher is more secure but slower
	bcryptCost = 12
)

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword checks if password meets requirements
func ValidatePassword(password string) bool {
	// Minimum 8 characters
	if len(password) < 8 {
		return false
	}
	return true
}
