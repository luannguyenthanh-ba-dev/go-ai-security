package usecase

import (
	"context"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/dto"
	usersRepository "github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/repository"
	userUseCase "github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/usecase"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
)

// Auth use case (application service)
type AuthService interface {
	Login(ctx context.Context, data *dto.LoginRequest) (*domain.JWTAuthEntity, error)
}

type authService struct {
	userService userUseCase.UserService
	jwtService  JWTService
}

func NewAuthService(userService userUseCase.UserService, jwtService JWTService) AuthService {
	return &authService{userService: userService, jwtService: jwtService}
}

func (service *authService) Login(ctx context.Context, data *dto.LoginRequest) (*domain.JWTAuthEntity, error) {
	user, err := service.userService.FindAUserByFilters(ctx, usersRepository.UserFilters{
		Username: &data.Username,
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrAuthUserNotFound
	}
	// Compare password
	if !utils.ComparePassword(data.Password, user.Password) {
		return nil, domain.ErrInvalidPassword
	}
	// Generate JWT
	auth, err := service.jwtService.GenerateJWT(&Claims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Phone:    user.Phone,
		Address:  user.Address,
		Gender:   user.Gender,
	})
	if err != nil {
		return nil, err
	}
	return auth, nil
}
