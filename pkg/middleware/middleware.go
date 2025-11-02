package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/usecase"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
	"go.uber.org/zap"
)

// HTTP middleware functions

var (
	ErrAuthorizationHeaderRequired = utils.NewCustomError(
		"AUTHORIZATION_HEADER_REQUIRED",
		http.StatusUnauthorized,
		"required authorization header",
	)
	ErrAuthorizationHeaderInvalid = utils.NewCustomError(
		"AUTHORIZATION_HEADER_INVALID",
		http.StatusUnauthorized,
		"invalid authorization header",
	)
	ErrAuthorizationInternalServerError = utils.NewCustomError(
		"AUTHORIZATION_INTERNAL_SERVER_ERROR",
		http.StatusInternalServerError,
		"internal server error while authorizing",
	)
)

type Middleware interface {
	AuthJWTMiddleware() gin.HandlerFunc
}

type middleware struct {
	jwtService usecase.JWTService
}

func NewMiddleware(jwtService usecase.JWTService) Middleware {
	return &middleware{jwtService: jwtService}
}

func (m *middleware) AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			zap.L().Error("Authorization header is required")
			utils.ErrorResponse(c,
				ErrAuthorizationHeaderRequired.HTTPStatus(),
				ErrAuthorizationHeaderRequired.Code(),
				ErrAuthorizationHeaderRequired.Error(),
			)
			c.Abort() // Stop processing the request chain
			return
		}
		if !strings.HasPrefix(token, "Bearer ") {
			zap.L().Error("Authorization header is invalid")
			utils.ErrorResponse(c,
				ErrAuthorizationHeaderInvalid.HTTPStatus(),
				ErrAuthorizationHeaderInvalid.Code(),
				ErrAuthorizationHeaderInvalid.Error(),
			)
			c.Abort() // Stop processing the request chain
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := m.jwtService.VerifyAccessToken(token)
		if err != nil {
			zap.L().Error("Failed to verify token", zap.Error(err))
			if ce, ok := err.(*utils.CustomError); ok {
				utils.ErrorResponse(c, ce.HTTPStatus(), ce.Code(), ce.Error())
				c.Abort() // Stop processing the request chain
				return
			}
			utils.ErrorResponse(c,
				ErrAuthorizationInternalServerError.HTTPStatus(),
				ErrAuthorizationInternalServerError.Code(),
				ErrAuthorizationInternalServerError.Error(),
			)
			c.Abort() // Stop processing the request chain
			return
		}

		c.Set("claims", claims)
		c.Set("request_user_id", claims.UserID)
		c.Next()
	}
}
