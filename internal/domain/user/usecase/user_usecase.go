package usecase

import (
	"context"
	"fmt"

	"github.com/TubagusAldiMY/go-template/internal/domain/user/dto"
	"github.com/TubagusAldiMY/go-template/internal/domain/user/entity"
	"github.com/TubagusAldiMY/go-template/internal/domain/user/repository"
	"github.com/TubagusAldiMY/go-template/internal/infrastructure/cache"
	"github.com/TubagusAldiMY/go-template/internal/shared/constants"
	"github.com/TubagusAldiMY/go-template/internal/shared/errors"
	"github.com/TubagusAldiMY/go-template/pkg/crypto"
	"github.com/TubagusAldiMY/go-template/pkg/jwt"
	"github.com/TubagusAldiMY/go-template/pkg/logger"
	"go.uber.org/zap"
)

type UserUsecase struct {
	userRepo       repository.UserRepository
	passwordHasher *crypto.PasswordHasher
	jwtManager     *jwt.Manager
	cache          *cache.Redis
}

func NewUserUsecase(
	userRepo repository.UserRepository,
	passwordHasher *crypto.PasswordHasher,
	jwtManager *jwt.Manager,
	cache *cache.Redis,
) *UserUsecase {
	return &UserUsecase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		jwtManager:     jwtManager,
		cache:          cache,
	}
}

func (uc *UserUsecase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("failed to check email existence", zap.Error(err))
		return nil, errors.ErrInternal
	}
	if exists {
		return nil, errors.ErrEmailAlreadyExists
	}

	// Check if username already exists
	exists, err = uc.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		logger.Error("failed to check username existence", zap.Error(err))
		return nil, errors.ErrInternal
	}
	if exists {
		return nil, errors.ErrUsernameAlreadyExists
	}

	// Hash password
	hashedPassword, err := uc.passwordHasher.Hash(req.Password)
	if err != nil {
		logger.Error("failed to hash password", zap.Error(err))
		return nil, errors.ErrInternal
	}

	// Create user entity
	user := entity.NewUser(req.Email, req.Username, hashedPassword, req.FullName, constants.RoleUser)

	// Save to database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		logger.Error("failed to create user", zap.Error(err))
		return nil, errors.ErrInternal
	}

	logger.Info("user registered successfully",
		zap.String("user_id", user.ID),
		zap.String("email", user.Email),
	)

	return uc.toUserResponse(user), nil
}

func (uc *UserUsecase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errors.ErrUserNotFound) {
			return nil, errors.ErrInvalidCredentials
		}
		logger.Error("failed to get user by email", zap.Error(err))
		return nil, errors.ErrInternal
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.ErrUnauthorized
	}

	// Verify password
	if !uc.passwordHasher.IsValid(user.Password, req.Password) {
		return nil, errors.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := uc.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		logger.Error("failed to generate access token", zap.Error(err))
		return nil, errors.ErrInternal
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		logger.Error("failed to generate refresh token", zap.Error(err))
		return nil, errors.ErrInternal
	}

	logger.Info("user logged in successfully",
		zap.String("user_id", user.ID),
		zap.String("email", user.Email),
	)

	return &dto.LoginResponse{
		User:         uc.toUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutes
	}, nil
}

func (uc *UserUsecase) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	// Validate refresh token
	userID, err := uc.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errors.ErrUserNotFound) {
			return nil, errors.ErrUnauthorized
		}
		logger.Error("failed to get user by id", zap.Error(err))
		return nil, errors.ErrInternal
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.ErrUnauthorized
	}

	// Generate new tokens
	accessToken, err := uc.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		logger.Error("failed to generate access token", zap.Error(err))
		return nil, errors.ErrInternal
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		logger.Error("failed to generate refresh token", zap.Error(err))
		return nil, errors.ErrInternal
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}, nil
}

func (uc *UserUsecase) GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("%s%s", constants.CacheKeyUserPrefix, userID)

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errors.ErrUserNotFound) {
			return nil, errors.ErrUserNotFound
		}
		logger.Error("failed to get user profile", zap.Error(err))
		return nil, errors.ErrInternal
	}

	// Cache the user
	_ = uc.cache.Set(ctx, cacheKey, user, constants.CacheTTLMedium)

	return uc.toUserResponse(user), nil
}

func (uc *UserUsecase) UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errors.ErrUserNotFound) {
			return nil, errors.ErrUserNotFound
		}
		logger.Error("failed to get user", zap.Error(err))
		return nil, errors.ErrInternal
	}

	user.UpdateProfile(req.FullName)

	if err := uc.userRepo.Update(ctx, user); err != nil {
		logger.Error("failed to update user", zap.Error(err))
		return nil, errors.ErrInternal
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("%s%s", constants.CacheKeyUserPrefix, userID)
	_ = uc.cache.Delete(ctx, cacheKey)

	logger.Info("user profile updated",
		zap.String("user_id", userID),
	)

	return uc.toUserResponse(user), nil
}

func (uc *UserUsecase) ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errors.ErrUserNotFound) {
			return errors.ErrUserNotFound
		}
		logger.Error("failed to get user", zap.Error(err))
		return errors.ErrInternal
	}

	// Verify old password
	if !uc.passwordHasher.IsValid(user.Password, req.OldPassword) {
		return errors.ErrInvalidPassword
	}

	// Hash new password
	hashedPassword, err := uc.passwordHasher.Hash(req.NewPassword)
	if err != nil {
		logger.Error("failed to hash password", zap.Error(err))
		return errors.ErrInternal
	}

	user.UpdatePassword(hashedPassword)

	if err := uc.userRepo.Update(ctx, user); err != nil {
		logger.Error("failed to update password", zap.Error(err))
		return errors.ErrInternal
	}

	logger.Info("password changed successfully",
		zap.String("user_id", userID),
	)

	return nil
}

func (uc *UserUsecase) ListUsers(ctx context.Context, req *dto.ListUsersRequest) ([]*dto.UserResponse, int64, error) {
	users, total, err := uc.userRepo.List(ctx, req.Page, req.PageSize, req.Search, req.Role, req.Status)
	if err != nil {
		logger.Error("failed to list users", zap.Error(err))
		return nil, 0, errors.ErrInternal
	}

	responses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = uc.toUserResponse(user)
	}

	return responses, total, nil
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, userID string) error {
	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		if errors.Is(err, errors.ErrUserNotFound) {
			return errors.ErrUserNotFound
		}
		logger.Error("failed to delete user", zap.Error(err))
		return errors.ErrInternal
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("%s%s", constants.CacheKeyUserPrefix, userID)
	_ = uc.cache.Delete(ctx, cacheKey)

	logger.Info("user deleted successfully",
		zap.String("user_id", userID),
	)

	return nil
}

func (uc *UserUsecase) toUserResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FullName:  user.FullName,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
