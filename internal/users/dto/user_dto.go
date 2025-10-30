package dto

import "github.com/luannguyenthanh-ba-dev/go-ai-security/internal/shared"

// User DTOs for request/response

type CreateUserRequest struct {
	Username string        `json:"username" binding:"required,min=5,max=20"`                     // required, min 5 characters, max 20 characters
	Email    string        `json:"email" binding:"required,email"`                               // required, email format
	Password string        `json:"password" binding:"required,min=6,max=20"`                     // required, min 6 characters, max 20 characters
	Name     string        `json:"name" binding:"required,min=3,max=50"`                         // required, min 3 characters, max 50 characters
	Phone    string        `json:"phone" binding:"omitempty,min=10,max=15" example:"0912345678"` // optional, min 10 characters, max 15 characters
	Address  string        `json:"address" binding:"omitempty,max=255"`                          // optional, max 255 characters
	Gender   shared.Gender `json:"gender" binding:"omitempty,oneof=1 2 3"`                       // optional, one of 1, 2, 3
}
