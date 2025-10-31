package usecase

import (
	"context"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/repository"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// User use case (application service)

type UserService interface {
	CreateUser(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error)
	FindAUserByFilters(ctx context.Context, filters repository.UserFilters) (*domain.UserEntity, error)
}

type userService struct {
	repo       repository.UserRepository
	saltRounds int
}

func NewUserService(r repository.UserRepository, saltRounds int) UserService {
	return &userService{repo: r, saltRounds: saltRounds}
}

func (service *userService) CreateUser(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error) {
	g, ctx := errgroup.WithContext(ctx)

	// Check existing user by username or email
	g.Go(func() error {
		existingEmailUser, err := service.repo.FindAUserByFilters(ctx, repository.UserFilters{
			Email: &user.Email,
		})
		if err != nil {
			return err
		}
		if existingEmailUser != nil {
			zap.L().Error("email already exists", zap.String("email", user.Email))
			return domain.ErrUserEmailAlreadyExists
		}
		return nil
	})

	// Check existing user by username
	g.Go(func() error {
		existingUsernameUser, err := service.repo.FindAUserByFilters(ctx, repository.UserFilters{
			Username: &user.Username,
		})
		if err != nil {
			return err
		}
		if existingUsernameUser != nil {
			zap.L().Error("username already exists", zap.String("username", user.Username))
			return domain.ErrUserUsernameAlreadyExists
		}
		return nil
	})

	// Wait for all checks to complete
	if err := g.Wait(); err != nil {
		zap.L().Error("error waiting for all checks to complete", zap.Error(err))
		return nil, err // It will return the first error that occurred and stop the execution
	}

	// Create user
	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password, service.saltRounds)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	user, err = service.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Clear password from response
	user.Password = ""

	return user, nil
}

func (service *userService) FindAUserByFilters(ctx context.Context, filters repository.UserFilters) (*domain.UserEntity, error) {
	user, err := service.repo.FindAUserByFilters(ctx, filters)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}
