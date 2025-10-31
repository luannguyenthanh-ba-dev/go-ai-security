package http

// HTTP routes configuration

import (
	"github.com/gin-gonic/gin"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/usecase"
)

func RegisterUserRoutes(router *gin.RouterGroup, userService usecase.UserService) {
	users := router.Group("/users")
	userHandler := NewUserHandler(userService)
	{
		users.POST("/register", userHandler.RegisterUser)
		// users.GET("/me", NewUserHandler(userService).GetMe)
		users.GET("/:id", userHandler.ViewUserInformation)
	}
}
