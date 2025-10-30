package utils

import "github.com/gin-gonic/gin"

// Response functions

func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, map[string]interface{}{
		"statusCode": statusCode,
		"message":    "success",
		"data":       data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, map[string]interface{}{
		"statusCode": statusCode,
		"error":      message,
		"code":       code,
	})
}
