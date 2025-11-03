package usecase

import (
	"context"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/dto"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/repository"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/shared"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// User use case (application service)

type UserService interface {
	RegisterUser(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error)
	FindAUserByFilters(ctx context.Context, filters repository.UserFilters) (*domain.UserEntity, error)
	UpdateUserProfile(ctx context.Context, userID *primitive.ObjectID, data *dto.UpdateMyProfileRequest) (*dto.UpdateMyProfileResponse, error)
}

type userService struct {
	repo       repository.UserRepository
	saltRounds int
}

func NewUserService(r repository.UserRepository, saltRounds int) UserService {
	return &userService{repo: r, saltRounds: saltRounds}
}

func (uSvc *userService) RegisterUser(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error) {
	// Use separate context for parallel checks
	// This allows canceling the checks if one fails, but preserves original context for CreateUser
	g, checkCtx := errgroup.WithContext(ctx)

	// Check existing user by username or email
	g.Go(func() error {
		// Create context with timeout for read operation within goroutine
		findCtx, findCancel := context.WithTimeout(checkCtx, shared.TimeoutForReadOperation)
		defer findCancel()

		existingEmailUser, err := uSvc.repo.FindOneByFilters(findCtx, repository.UserFilters{
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
		// Create context with timeout for read operation within goroutine
		findCtx, findCancel := context.WithTimeout(checkCtx, shared.TimeoutForReadOperation)
		defer findCancel()

		existingUsernameUser, err := uSvc.repo.FindOneByFilters(findCtx, repository.UserFilters{
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
	hashedPassword, err := utils.HashPassword(user.Password, uSvc.saltRounds)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	// Create context with timeout for write operation
	createCtx, createCancel := context.WithTimeout(ctx, shared.TimeoutForWriteOperation)
	defer createCancel()

	// Use context with timeout (not the errgroup context) for CreateUser
	// The errgroup context may be canceled after g.Wait(), but we still need to create the user
	user, err = uSvc.repo.CreateNew(createCtx, user)
	if err != nil {
		return nil, err
	}

	// Clear password from response
	user.Password = ""

	return user, nil
}

func (uSvc *userService) FindAUserByFilters(ctx context.Context, filters repository.UserFilters) (*domain.UserEntity, error) {
	// Create context with timeout for read operation
	findCtx, findCancel := context.WithTimeout(ctx, shared.TimeoutForReadOperation)
	defer findCancel()

	user, err := uSvc.repo.FindOneByFilters(findCtx, filters)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (uSvc *userService) UpdateUserProfile(
	ctx context.Context,
	userID *primitive.ObjectID,
	data *dto.UpdateMyProfileRequest) (*dto.UpdateMyProfileResponse, error) {
	// Create context with timeout for complex operation (read + write)
	updateCtx, updateCancel := context.WithTimeout(ctx, shared.TimeoutForComplexOperation)
	defer updateCancel()

	// Find user first
	user, err := uSvc.repo.FindOneByFilters(updateCtx, repository.UserFilters{ID: userID})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	// Update user profile
	updateData := &domain.UserEntity{
		Name:    data.Name,
		Phone:   data.Phone,
		Address: data.Address,
		Gender:  data.Gender,
		Avatar:  data.Avatar,
	}
	ok, err := uSvc.repo.UpdateBasicInfoByID(updateCtx, userID, updateData)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrUserUpdateFailed
	}
	return &dto.UpdateMyProfileResponse{Updated: true}, nil
}
