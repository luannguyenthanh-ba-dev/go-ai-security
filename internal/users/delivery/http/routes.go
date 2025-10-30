package http

// HTTP routes configuration

import (
	"github.com/gin-gonic/gin"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/usecase"
)

func RegisterUserRoutes(router *gin.RouterGroup, userService usecase.UserService) {
	users := router.Group("/users")
	{
		users.POST("/register", NewUserHandler(userService).RegisterUser)
	}
}
