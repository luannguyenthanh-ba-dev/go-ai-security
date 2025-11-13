package usecase

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/repository"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/shared"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
)

// JWT

type RefreshClaims struct {
	UserID string `json:"user_id" required:"true"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateJWT(claims *shared.Claims) (*domain.JWTAuthEntity, error)
	VerifyAccessToken(token string) (*shared.Claims, error)
	AddAccessTokenToWhiteList(ctx context.Context, userID string, accessToken string, ttl time.Duration) (bool, error)
	VerifyAccessTokenInWhiteList(ctx context.Context, userID string, token string) (bool, error)
	VerifyRefreshToken(refreshToken string) (*RefreshClaims, error)
}

type jwtService struct {
	secret              string
	expiresIn           time.Duration
	authCacheRepository repository.AuthCacheRepository
}

func NewJWTService(secret string, expiresIn time.Duration, authCacheRepository repository.AuthCacheRepository) JWTService {
	return &jwtService{
		secret:              secret,
		expiresIn:           expiresIn,
		authCacheRepository: authCacheRepository,
	}
}

func (jService *jwtService) GenerateJWT(claims *shared.Claims) (*domain.JWTAuthEntity, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(jService.expiresIn))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	// Generate access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(jService.secret))
	if err != nil {
		return nil, domain.ErrSigningAccessTokenFailed
	}

	// Generate refresh token
	refreshTokenClaims := &RefreshClaims{
		UserID: claims.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jService.secret))
	if err != nil {
		return nil, domain.ErrSigningRefreshTokenFailed
	}

	return &domain.JWTAuthEntity{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiredIn:    int64(jService.expiresIn.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (jService *jwtService) VerifyAccessToken(accessToken string) (*shared.Claims, error) {
	// Parse the token with the secret key
	token, err := jwt.ParseWithClaims(accessToken, &shared.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method to prevent algorithm confusion attacks
		// Only accept HS256 (same as when generating tokens)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrJWTTokenInvalidSigningMethod
		}
		return []byte(jService.secret), nil
	})

	// Handle parsing errors
	if err != nil {
		// Check if error is due to token expiration
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrJWTTokenExpired
		}
		// Other parsing errors (invalid signature, malformed token, etc.)
		return nil, domain.ErrorParsingJWTToken
	}

	// Verify token is valid (signature verified and not expired)
	if !token.Valid {
		return nil, domain.ErrJWTTokenInvalid
	}

	// Type assertion to extract custom claims
	claims, ok := token.Claims.(*shared.Claims)
	if !ok {
		return nil, domain.ErrInvalidJWTTokenClaims
	}

	// Additional expiration check (defensive programming)
	// jwt/v5 automatically validates expiration, but we check explicitly for clarity
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, domain.ErrJWTTokenExpired
	}

	return claims, nil
}

func (jService *jwtService) AddAccessTokenToWhiteList(ctx context.Context, userID string, accessToken string, ttl time.Duration) (bool, error) {
	return jService.authCacheRepository.AddAccessTokenToWhiteList(ctx, userID, accessToken, ttl)
}

func (jService *jwtService) VerifyAccessTokenInWhiteList(ctx context.Context, userID string, token string) (bool, error) {
	tokenInWhiteList, err := jService.authCacheRepository.GetAccessTokenFromWhiteList(ctx, userID)
	if err != nil {
		return false, utils.NewCustomError(
			"ERROR_VERIFY_ACCESS_TOKEN_IN_WHITE_LIST",
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	if tokenInWhiteList == "" {
		return false, domain.ErrTokenNotInWhiteList
	}
	if tokenInWhiteList != token {
		return false, domain.ErrNotMatchTokenInWhiteList
	}
	return true, nil
}

func (jService *jwtService) VerifyRefreshToken(refreshToken string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrJWTRefreshTokenInvalidSigningMethod
		}
		return []byte(jService.secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrJWTRefreshTokenExpired
		}
		return nil, domain.ErrorParsingJWTRefreshToken
	}
	if !token.Valid {
		return nil, domain.ErrJWTRefreshTokenInvalid
	}

	// Type assertion to extract custom claims
	claims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		return nil, domain.ErrInvalidJWTRefreshTokenClaims
	}
	// Additional expiration check (defensive programming)
	// jwt/v5 automatically validates expiration, but we check explicitly for clarity
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, domain.ErrJWTRefreshTokenExpired
	}
	return claims, nil
}
