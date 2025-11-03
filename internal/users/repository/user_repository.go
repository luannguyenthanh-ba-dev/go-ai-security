package repository

import (
	"context"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User repository interface

type UserRepository interface {
	CreateNew(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error)
	FindOneByFilters(ctx context.Context, filters UserFilters) (*domain.UserEntity, error)
	UpdateBasicInfoByID(ctx context.Context, userID *primitive.ObjectID, updateData *domain.UserEntity) (bool, error)
}
