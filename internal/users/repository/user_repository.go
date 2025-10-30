package repository

import (
	"context"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/domain"
)

// User repository interface

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error)
	FindAUserByFilters(ctx context.Context, filters UserFilters) (*domain.UserEntity, error)
}
