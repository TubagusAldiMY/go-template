package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TubagusAldiMY/go-template/internal/domain/user/dto"
	"github.com/TubagusAldiMY/go-template/internal/domain/user/entity"
	"github.com/TubagusAldiMY/go-template/internal/domain/user/usecase"
	sharedErrors "github.com/TubagusAldiMY/go-template/internal/shared/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, page, pageSize int, search, role, status string) ([]*entity.User, int64, error) {
	args := m.Called(ctx, page, pageSize, search, role, status)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

// MockPasswordHasher is a mock implementation of PasswordHasher
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Compare(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockPasswordHasher) IsValid(hashedPassword, password string) bool {
	args := m.Called(hashedPassword, password)
	return args.Bool(0)
}

// MockJWTManager is a mock implementation of JWTManager
type MockJWTManager struct {
	mock.Mock
}

func (m *MockJWTManager) GenerateAccessToken(userID, email, role string) (string, error) {
	args := m.Called(userID, email, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) GenerateRefreshToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

// MockRedis is a mock implementation of Redis
type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedis) Delete(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func TestRegister_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	mockJWT := new(MockJWTManager)
	mockRedis := new(MockRedis)

	uc := usecase.NewUserUsecase(mockRepo, mockHasher, mockJWT, mockRedis)

	req := &dto.RegisterRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "SecurePass123!",
		FullName: "Test User",
	}

	mockRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(false, nil)
	mockRepo.On("ExistsByUsername", mock.Anything, req.Username).Return(false, nil)
	mockHasher.On("Hash", req.Password).Return("hashedpassword", nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	// Act
	result, err := uc.Register(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Email, result.Email)
	assert.Equal(t, req.Username, result.Username)
	assert.Equal(t, req.FullName, result.FullName)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	mockJWT := new(MockJWTManager)
	mockRedis := new(MockRedis)

	uc := usecase.NewUserUsecase(mockRepo, mockHasher, mockJWT, mockRedis)

	req := &dto.RegisterRequest{
		Email:    "existing@example.com",
		Username: "testuser",
		Password: "SecurePass123!",
		FullName: "Test User",
	}

	mockRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(true, nil)

	// Act
	result, err := uc.Register(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, sharedErrors.ErrEmailAlreadyExists))

	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	mockJWT := new(MockJWTManager)
	mockRedis := new(MockRedis)

	uc := usecase.NewUserUsecase(mockRepo, mockHasher, mockJWT, mockRedis)

	req := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	user := &entity.User{
		ID:       "user-123",
		Email:    req.Email,
		Username: "testuser",
		Password: "hashedpassword",
		FullName: "Test User",
		Role:     "user",
		Status:   "active",
	}

	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(user, nil)
	mockHasher.On("IsValid", user.Password, req.Password).Return(true)
	mockJWT.On("GenerateAccessToken", user.ID, user.Email, user.Role).Return("access-token", nil)
	mockJWT.On("GenerateRefreshToken", user.ID).Return("refresh-token", nil)

	// Act
	result, err := uc.Login(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Equal(t, "refresh-token", result.RefreshToken)
	assert.Equal(t, "Bearer", result.TokenType)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	mockJWT := new(MockJWTManager)
	mockRedis := new(MockRedis)

	uc := usecase.NewUserUsecase(mockRepo, mockHasher, mockJWT, mockRedis)

	req := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "WrongPassword",
	}

	user := &entity.User{
		ID:       "user-123",
		Email:    req.Email,
		Password: "hashedpassword",
		Status:   "active",
	}

	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(user, nil)
	mockHasher.On("IsValid", user.Password, req.Password).Return(false)

	// Act
	result, err := uc.Login(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, sharedErrors.ErrInvalidCredentials))

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}
