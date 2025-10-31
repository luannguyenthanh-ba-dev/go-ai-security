package dto

// Auth DTOs for request/response
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=5,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

