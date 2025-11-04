package http

// HTTP routes configuration

import (
	"github.com/gin-gonic/gin"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/usecase"
)

func RegisterUserPublicRoutes(router *gin.RouterGroup, userService usecase.UserService) {
	userHandler := NewUserHandler(userService)
	users := router.Group("/users")
	{
		users.POST("/register", userHandler.RegisterUser)
		users.GET("/:id", userHandler.ViewUserInformation)
	}
}

func RegisterUserProtectedRoutes(router *gin.RouterGroup, userService usecase.UserService) {
	userHandler := NewUserHandler(userService)
	users := router.Group("/users")
	{
		users.GET("/me", userHandler.GetMe)
		users.PUT("/me", userHandler.UpdateMe)
		users.PUT("/me/passwords", userHandler.UpdateMyPassword)
	}
}
