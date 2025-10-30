package domain

import "errors"

// Domain-specific errors

var (
	// Invalid input errors
	ErrUserInvalidUsername = errors.New("invalid username")
	ErrUserInvalidEmail    = errors.New("invalid email")
	ErrUserInvalidPassword = errors.New("invalid password")
	ErrUserInvalidName     = errors.New("invalid name")
	ErrUserInvalidPhone    = errors.New("invalid phone")
	ErrUserInvalidGender   = errors.New("invalid gender")
	ErrUserInvalidInput    = errors.New("invalid input data")

	// Conflict errors
	ErrUserUsernameAlreadyExists = errors.New("username already exists")
	ErrUserEmailAlreadyExists    = errors.New("email already exists")
	ErrUserWrongPassword         = errors.New("wrong password")
	ErrUserNotHasRole            = errors.New("user does not have the required role")

	// Not found errors
	ErrUserNotFound = errors.New("user not found")

	// Internal server errors
	ErrUserInternalServerError = errors.New("internal server error")
)
