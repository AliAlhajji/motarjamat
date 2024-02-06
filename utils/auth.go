package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Generate a random salt using bcrypt's built-in functionality
func HashAndSaltPassword(password string) (string, error) {
	salt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(salt), nil
}

// Compare the password with the hashed password (including the salt)
func VerifyPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
