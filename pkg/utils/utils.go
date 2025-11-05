package utils

// Utility functions

import (
	"errors"
	"math/rand"

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

func RandomInt64(min, max int64) (int64, error) {
	if min > max {
		return 0, errors.New("min is greater than max")
	}
	return rand.Int63n(max-min+1) + min, nil
}
