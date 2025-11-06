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
	ErrInvalidJWTTokenClaims = utils.NewCustomError("INVALID_JWT_TOKEN_CLAIMS",
		http.StatusUnauthorized,
		"invalid jwt token claims",
	)
	ErrJWTTokenInvalidSigningMethod = utils.NewCustomError("JWT_TOKEN_INVALID_SIGNING_METHOD",
		http.StatusUnauthorized,
		"invalid jwt token signing method",
	)
	ErrorParsingJWTToken = utils.NewCustomError("ERROR_PARSING_JWT_TOKEN",
		http.StatusUnauthorized,
		"error parsing jwt token",
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

	// Forbidden errors
	ErrTokenNotInWhiteList = utils.NewCustomError("TOKEN_NOT_IN_WHITE_LIST", http.StatusForbidden, "token not in white list")
	ErrNotMatchTokenInWhiteList = utils.NewCustomError("NOT_MATCH_TOKEN_IN_WHITE_LIST", http.StatusForbidden, "not match token in white list")
)
