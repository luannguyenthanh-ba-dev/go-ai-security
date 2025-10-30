package http

// HTTP handlers for user endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/dto"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/usecase"
)

// UserHandler holds dependencies for user HTTP handlers
// e.g., the usecase layer

type UserHandler struct {
	usecase usecase.UserUseCase
}

func NewUserHandler(usecase usecase.UserUseCase) *UserHandler {
	return &UserHandler{usecase: usecase}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	user, err := h.usecase.CreateUser(c.Request.Context(), userEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
