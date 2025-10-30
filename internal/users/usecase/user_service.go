package usecase

import (
	"context"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/repository"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
)

// User use case (application service)

type UserService interface {
	CreateUser(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error)
}

type userService struct {
	repo       repository.UserRepository
	saltRounds int
}

func NewUserService(r repository.UserRepository, saltRounds int) UserService {
	return &userService{repo: r, saltRounds: saltRounds}
}

func (u *userService) CreateUser(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error) {
	// Check existing user by username or email

	// Create user
	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password, u.saltRounds)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	user, err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Clear password from response
	user.Password = ""

	return user, nil
}
