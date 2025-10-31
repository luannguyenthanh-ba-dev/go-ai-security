package domain

// User domain entity

import (
	"time"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserEntity struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string             `bson:"username,required" json:"username" binding:"required,min=5,max=20"`
	Email     string             `bson:"email,required" json:"email" binding:"required,email"`
	Password  string             `bson:"password,required" json:"-"`
	Name      string             `bson:"name,required" json:"name" binding:"required,min=3,max=50"`
	Phone     string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Address   string             `bson:"address,omitempty" json:"address,omitempty"`
	Role      shared.Role        `bson:"role,required" json:"role" binding:"required"`
	Gender    shared.Gender      `bson:"gender,omitempty" json:"gender,omitempty"`
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
