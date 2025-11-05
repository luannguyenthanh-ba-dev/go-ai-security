package usecase

import (
	"context"
	"time"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/dto"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/repository"
	usersRepository "github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/repository"
	userUseCase "github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/usecase"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/shared"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
	"go.uber.org/zap"
)

// Auth use case (application service)
type AuthService interface {
	Login(ctx context.Context, data *dto.LoginRequest) (*domain.JWTAuthEntity, error)
}

type authService struct {
	userService         userUseCase.UserService
	jwtService          JWTService
	authCacheRepository repository.AuthCacheRepository
}

func NewAuthService(
	userService userUseCase.UserService,
	jwtService JWTService,
	authCacheRepository repository.AuthCacheRepository) AuthService {
	return &authService{userService: userService, jwtService: jwtService, authCacheRepository: authCacheRepository}
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
	auth, err := service.jwtService.GenerateJWT(&shared.Claims{
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

	// Cache the access token in white list asynchronously
	// This allows cache operation to complete even after response is sent
	go func() {
		// Create background context for async operation
		// This prevents the goroutine from being cancelled when the original request context is done
		backgroundCtx := context.Background()

		// Generate a random TTL between 5 and 30 seconds
		randomTTL, _ := utils.RandomInt64(5, 30)
		ttl := time.Duration(auth.ExpiredIn) * time.Second + time.Duration(randomTTL) * time.Second

		// Cache the token
		success, err := service.authCacheRepository.AddAccessTokenToWhiteList(backgroundCtx, user.ID.Hex(), auth.AccessToken, ttl)

		// Log the result (non-blocking, doesn't affect response)
		if err != nil {
			zap.L().Error("failed to cache access token in whitelist",
				zap.String("userID", user.ID.Hex()),
				zap.Error(err),
			)
		} else if success {
			zap.L().Debug("access token cached in whitelist successfully",
				zap.String("userID", user.ID.Hex()),
				zap.Duration("ttl", ttl),
			)
		}
	}()

	return auth, nil
}
