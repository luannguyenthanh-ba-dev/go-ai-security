package domain

import (
	"net/http"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
)

// Domain-specific errors

var (
	// Invalid input errors
	ErrUserInvalidUsername = utils.NewCustomError("USER_INVALID_USERNAME", http.StatusBadRequest, "invalid username")
	ErrUserInvalidEmail    = utils.NewCustomError("USER_INVALID_EMAIL", http.StatusBadRequest, "invalid email")
	ErrUserInvalidPassword = utils.NewCustomError("USER_INVALID_PASSWORD", http.StatusBadRequest, "invalid password")
	ErrUserInvalidName     = utils.NewCustomError("USER_INVALID_NAME", http.StatusBadRequest, "invalid name")
	ErrUserInvalidPhone    = utils.NewCustomError("USER_INVALID_PHONE", http.StatusBadRequest, "invalid phone")
	ErrUserInvalidGender   = utils.NewCustomError("USER_INVALID_GENDER", http.StatusBadRequest, "invalid gender")
	ErrUserInvalidInput    = utils.NewCustomError("USER_INVALID_INPUT", http.StatusBadRequest, "invalid input data")
	ErrUserInvalidID       = utils.NewCustomError("USER_INVALID_ID", http.StatusBadRequest, "invalid user id")

	// Conflict errors
	ErrUserUsernameAlreadyExists = utils.NewCustomError("USER_USERNAME_ALREADY_EXISTS", http.StatusConflict, "username already exists")
	ErrUserEmailAlreadyExists    = utils.NewCustomError("USER_EMAIL_ALREADY_EXISTS", http.StatusConflict, "email already exists")
	ErrUserWrongPassword         = utils.NewCustomError("USER_WRONG_PASSWORD", http.StatusUnauthorized, "wrong password")
	ErrUserNotHasRole            = utils.NewCustomError("USER_NOT_HAS_ROLE", http.StatusForbidden, "user does not have the required role")

	// Not found errors
	ErrUserNotFound = utils.NewCustomError("USER_NOT_FOUND", http.StatusNotFound, "user not found")

	// Internal server errors
	ErrUserInternalServerError = utils.NewCustomError("USER_INTERNAL_SERVER_ERROR", http.StatusInternalServerError, "internal server error")
)
