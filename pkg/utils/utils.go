package utils

// Utility functions

import (
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string, saltRounds int) (string, error) {
	// Check if the password is empty
	if password == "" {
		return "", errors.New("password is empty")
	}

	// Check if the salt rounds is valid
	if saltRounds <= 0 {
		return "", errors.New("invalid salt rounds")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), saltRounds)
	if err != nil {
		zap.L().Error("error hashing password", zap.Error(err))
		return "", errors.New("error hashing password")
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a password with a hashed password
func ComparePassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
