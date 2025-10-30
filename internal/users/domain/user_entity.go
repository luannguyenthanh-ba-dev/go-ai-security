package domain

// User domain entity

import (
	"time"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserEntity struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string             `bson:"username,omitempty" json:"username,omitempty"`
	Email     string             `bson:"email,omitempty" json:"email,omitempty"`
	Password  string             `bson:"password,omitempty" json:"-"` // Password is not returned in the response
	Name      string             `bson:"name,omitempty" json:"name,omitempty"`
	Phone     string             `bson:"phone" json:"phone"`
	Address   string             `bson:"address" json:"address"`
	Role      shared.Role        `bson:"role" json:"role"`
	Gender    shared.Gender      `bson:"gender" json:"gender"`
	CreatedAt int64              `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt int64              `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// NewUserEntity is a constructor for the UserEntity struct
func NewUserEntity(username, email, password, name, phone, address string, role shared.Role, gender shared.Gender) *UserEntity {
	return &UserEntity{
		Username:  username,
		Email:     email,
		Password:  password,
		Name:      name,
		Phone:     phone,
		Address:   address,
		Role:      role,
		Gender:    gender,
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
}
