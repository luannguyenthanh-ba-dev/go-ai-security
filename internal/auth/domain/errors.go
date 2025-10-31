package domain

import (
	"net/http"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
)

// Domain-specific errors
var (
	ErrJWTTokenInvalid = utils.NewCustomError("JWT_TOKEN_INVALID",
		http.StatusUnauthorized,
		"invalid jwt token",
	)
	ErrJWTTokenExpired = utils.NewCustomError("JWT_TOKEN_EXPIRED",
		http.StatusUnauthorized,
		"jwt token expired",
	)
	ErrJWTRefreshTokenInvalid = utils.NewCustomError("JWT_REFRESH_TOKEN_INVALID",
		http.StatusUnauthorized,
		"invalid jwt refresh token",
	)
	ErrJWTRefreshTokenExpired = utils.NewCustomError("JWT_REFRESH_TOKEN_EXPIRED",
		http.StatusUnauthorized,
		"jwt refresh token expired",
	)

	// Not found errors
	ErrAuthUserNotFound = utils.NewCustomError("AUTH_USER_NOT_FOUND", http.StatusNotFound, "user not found")

	// Internal server errors
	ErrAuthInternalServerError = utils.NewCustomError("AUTH_INTERNAL_SERVER_ERROR",
		http.StatusInternalServerError,
		"internal server error",
	)
	ErrSigningAccessTokenFailed  = utils.NewCustomError("SIGNING_ACCESS_TOKEN_FAILED",
		http.StatusInternalServerError,
		"failed to sign access token",
	)
	ErrSigningRefreshTokenFailed = utils.NewCustomError("SIGNING_REFRESH_TOKEN_FAILED",
		http.StatusInternalServerError,
		"failed to sign refresh token",
	)

	// Unauthorized errors
	ErrInvalidPassword = utils.NewCustomError("INVALID_PASSWORD", http.StatusUnauthorized, "invalid password")
)
