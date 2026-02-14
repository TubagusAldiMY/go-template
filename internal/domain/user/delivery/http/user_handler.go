package http

import (
	"net/http"

	"github.com/TubagusAldiMY/go-template/internal/domain/user/dto"
	"github.com/TubagusAldiMY/go-template/internal/domain/user/usecase"
	"github.com/TubagusAldiMY/go-template/internal/shared/constants"
	"github.com/TubagusAldiMY/go-template/internal/shared/errors"
	"github.com/TubagusAldiMY/go-template/pkg/logger"
	"github.com/TubagusAldiMY/go-template/pkg/response"
	customValidator "github.com/TubagusAldiMY/go-template/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := customValidator.Validate(&req); err != nil {
		validationErrors := customValidator.FormatValidationErrors(err)
		response.UnprocessableEntity(c, "Validation failed", validationErrors)
		return
	}

	user, err := h.userUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrEmailAlreadyExists):
			response.Conflict(c, "Email already exists", nil)
		case errors.Is(err, errors.ErrUsernameAlreadyExists):
			response.Conflict(c, "Username already exists", nil)
		default:
			logger.Error("failed to register user", zap.Error(err))
			response.InternalServerError(c, "Failed to register user")
		}
		return
	}

	response.Created(c, "User registered successfully", user)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and get tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := customValidator.Validate(&req); err != nil {
		validationErrors := customValidator.FormatValidationErrors(err)
		response.UnprocessableEntity(c, "Validation failed", validationErrors)
		return
	}

	loginResp, err := h.userUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrInvalidCredentials):
			response.Unauthorized(c, "Invalid email or password")
		case errors.Is(err, errors.ErrUnauthorized):
			response.Unauthorized(c, "Account is not active")
		default:
			logger.Error("failed to login", zap.Error(err))
			response.InternalServerError(c, "Failed to login")
		}
		return
	}

	response.OK(c, "Login successful", loginResp)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.Response{data=dto.RefreshTokenResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := customValidator.Validate(&req); err != nil {
		validationErrors := customValidator.FormatValidationErrors(err)
		response.UnprocessableEntity(c, "Validation failed", validationErrors)
		return
	}

	refreshResp, err := h.userUsecase.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrInvalidToken):
			response.Unauthorized(c, "Invalid refresh token")
		case errors.Is(err, errors.ErrUnauthorized):
			response.Unauthorized(c, "Unauthorized")
		default:
			logger.Error("failed to refresh token", zap.Error(err))
			response.InternalServerError(c, "Failed to refresh token")
		}
		return
	}

	response.OK(c, "Token refreshed successfully", refreshResp)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString(constants.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	user, err := h.userUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrUserNotFound):
			response.NotFound(c, "User not found")
		default:
			logger.Error("failed to get profile", zap.Error(err))
			response.InternalServerError(c, "Failed to get profile")
		}
		return
	}

	response.OK(c, "Profile retrieved successfully", user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString(constants.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := customValidator.Validate(&req); err != nil {
		validationErrors := customValidator.FormatValidationErrors(err)
		response.UnprocessableEntity(c, "Validation failed", validationErrors)
		return
	}

	user, err := h.userUsecase.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrUserNotFound):
			response.NotFound(c, "User not found")
		default:
			logger.Error("failed to update profile", zap.Error(err))
			response.InternalServerError(c, "Failed to update profile")
		}
		return
	}

	response.OK(c, "Profile updated successfully", user)
}

// ChangePassword godoc
// @Summary Change password
// @Description Change authenticated user's password
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.ChangePasswordRequest true "Change password request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString(constants.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := customValidator.Validate(&req); err != nil {
		validationErrors := customValidator.FormatValidationErrors(err)
		response.UnprocessableEntity(c, "Validation failed", validationErrors)
		return
	}

	if err := h.userUsecase.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		switch {
		case errors.Is(err, errors.ErrUserNotFound):
			response.NotFound(c, "User not found")
		case errors.Is(err, errors.ErrInvalidPassword):
			response.BadRequest(c, "Invalid old password", nil)
		default:
			logger.Error("failed to change password", zap.Error(err))
			response.InternalServerError(c, "Failed to change password")
		}
		return
	}

	response.OK(c, "Password changed successfully", nil)
}

// ListUsers godoc
// @Summary List users
// @Description Get list of users with pagination and filters (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search by email, username, or full name"
// @Param role query string false "Filter by role"
// @Param status query string false "Filter by status"
// @Success 200 {object} response.Response{data=[]dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "Invalid query parameters", err.Error())
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	if err := customValidator.Validate(&req); err != nil {
		validationErrors := customValidator.FormatValidationErrors(err)
		response.UnprocessableEntity(c, "Validation failed", validationErrors)
		return
	}

	users, total, err := h.userUsecase.ListUsers(c.Request.Context(), &req)
	if err != nil {
		logger.Error("failed to list users", zap.Error(err))
		response.InternalServerError(c, "Failed to list users")
		return
	}

	meta := response.NewMeta(req.Page, req.PageSize, total)
	response.SuccessWithMeta(c, "Users retrieved successfully", users, meta)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	if err := h.userUsecase.DeleteUser(c.Request.Context(), userID); err != nil {
		switch {
		case errors.Is(err, errors.ErrUserNotFound):
			response.NotFound(c, "User not found")
		default:
			logger.Error("failed to delete user", zap.Error(err))
			response.InternalServerError(c, "Failed to delete user")
		}
		return
	}

	response.OK(c, "User deleted successfully", nil)
}
