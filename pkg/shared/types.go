package shared

// Shared types and constants

import (
	"github.com/golang-jwt/jwt/v5"
)

// Shared roles
type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// Write a function to check if the role is valid
// IsValid is a method of the Role type that returns a boolean to check if the role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RoleUser:
		return true
	default:
		return false
	}
}

// Shared genders
type Gender int

const (
	GenderMale Gender = 1
	GenderFemale Gender = 2
	GenderUnknown Gender = 3
)

// Write a function to check if the gender is valid
// IsValid is a method of the Gender type that returns a boolean to check if the gender is valid
func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderUnknown:
		return true
	default:
		return false
	}
}

// Shared claims for JWT
type Claims struct {
	UserID   string        `json:"user_id" required:"true"`
	Username string        `json:"username" required:"true"`
	Email    string        `json:"email" required:"true"`
	Role     Role   `json:"role" required:"true"`
	Phone    string        `json:"phone,omitempty"`
	Address  string        `json:"address,omitempty"`
	Gender   Gender `json:"gender,omitempty"`
	jwt.RegisteredClaims
}


