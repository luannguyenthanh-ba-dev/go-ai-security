package usecase

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/shared"
)

// JWT

type Claims struct {
	UserID   string        `json:"user_id" required:"true"`
	Username string        `json:"username" required:"true"`
	Email    string        `json:"email" required:"true"`
	Role     shared.Role   `json:"role" required:"true"`
	Phone    string        `json:"phone,omitempty"`
	Address  string        `json:"address,omitempty"`
	Gender   shared.Gender `json:"gender,omitempty"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id" required:"true"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateJWT(claims *Claims) (*domain.JWTAuthEntity, error)
}

type jwtService struct {
	secret    string
	expiresIn time.Duration
}

func NewJWTService(secret string, expiresIn time.Duration) JWTService {
	return &jwtService{
		secret:    secret,
		expiresIn: expiresIn,
	}
}

func (jService *jwtService) GenerateJWT(claims *Claims) (*domain.JWTAuthEntity, error) {
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
