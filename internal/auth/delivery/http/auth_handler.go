package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/dto"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/usecase"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
	"go.uber.org/zap"
)

// HTTP handlers for auth endpoints

type  AuthHandler struct {
	service usecase.AuthService
}

func NewAuthHandler(service usecase.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login handles POST /auth/login request
// @Summary Login
// @Description Login with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Login request"
// @Success 201 {object} domain.JWTAuthEntity
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var data dto.LoginRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "AUTH_INVALID_INPUT", err.Error())
		return
	}

	zap.L().Info("Login request received", zap.Any("data", data))

	auth, err := h.service.Login(c.Request.Context(), &data)
	if err != nil {
		if ce, ok := err.(*utils.CustomError); ok {
			utils.ErrorResponse(c, ce.HTTPStatus(), ce.Code(), ce.Error())
			return
		}
		utils.ErrorResponse(c,
			domain.ErrAuthInternalServerError.HTTPStatus(),
			domain.ErrAuthInternalServerError.Code(),
			domain.ErrAuthInternalServerError.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, auth)
}
