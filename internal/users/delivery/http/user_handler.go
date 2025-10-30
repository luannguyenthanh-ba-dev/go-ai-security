package http

// HTTP handlers for user endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/dto"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/usecase"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/utils"
)

// UserHandler holds dependencies for user HTTP handlers
// e.g., the usecase layer

type UserHandler struct {
	service usecase.UserService
}

func NewUserHandler(service usecase.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterUser handles POST /users/register request
// @Summary Register a new user
// @Description Register a new user with the given information
// @Tags Users
// @Accept json
// @Produce json
// @Param body body dto.CreateUserRequest true "User information"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	// Using dto here
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "USER_INVALID_INPUT", err.Error())
		return
	}

	// Map DTO to domain entity
	userEntity := &domain.UserEntity{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
		Address:  req.Address,
		Gender:   req.Gender, // already shared.Gender
	}

	user, err := h.service.CreateUser(c.Request.Context(), userEntity)
	if err != nil {
		if ce, ok := err.(*utils.CustomError); ok {
			utils.ErrorResponse(c, ce.HTTPStatus(), ce.Code(), ce.Error())
			return
		}
		utils.ErrorResponse(c,
			domain.ErrUserInternalServerError.HTTPStatus(),
			domain.ErrUserInternalServerError.Code(),
			domain.ErrUserInternalServerError.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, user)
}
