package domain

// Auth domain entity
type JWTAuthEntity struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
	ExpiredIn    int64  `json:"expired_in" binding:"required"`
	TokenType    string `json:"token_type" binding:"required"`
}
