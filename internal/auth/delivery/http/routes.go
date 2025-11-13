package http

import (
	"github.com/gin-gonic/gin"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/usecase"
)

// HTTP routes configuration
func RegisterAuthPublicRoutes(router *gin.RouterGroup, authService usecase.AuthService) {
	authHandler := NewAuthHandler(authService)
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/tokens/refresh", authHandler.RefreshAccessToken)
	}
}

func RegisterAuthProtectedRoutes(router *gin.RouterGroup, authService usecase.AuthService) {
	// authHandler := NewAuthHandler(authService)
	// auth := router.Group("/auth")
	// {}
}
