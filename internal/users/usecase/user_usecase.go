package usecase

import (
	"context"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/repository"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
)

// User use case (application service)

type UserUseCase interface {
	CreateUser(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error)
}

type userUseCase struct {
	repo       repository.UserRepository
	saltRounds int
}

func NewUserUseCase(r repository.UserRepository, saltRounds int) UserUseCase {
	return &userUseCase{repo: r, saltRounds: saltRounds}
}

func (u *userUseCase) CreateUser(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error) {
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
