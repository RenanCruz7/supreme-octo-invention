package mocks

import (
	"awesomeProject/models"

	"github.com/stretchr/testify/mock"
)

// ────────────────────────────────────────────────────────────────
// MockAuthService  (implements services.IAuthService)
// ────────────────────────────────────────────────────────────────

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	args := m.Called(req)
	result, _ := args.Get(0).(*models.AuthResponse)
	return result, args.Error(1)
}

func (m *MockAuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	args := m.Called(req)
	result, _ := args.Get(0).(*models.AuthResponse)
	return result, args.Error(1)
}

func (m *MockAuthService) RefreshTokens(refreshToken string) (*models.AuthResponse, error) {
	args := m.Called(refreshToken)
	result, _ := args.Get(0).(*models.AuthResponse)
	return result, args.Error(1)
}

func (m *MockAuthService) Logout(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

// ────────────────────────────────────────────────────────────────
// MockUserService  (implements services.IUserService)
// ────────────────────────────────────────────────────────────────

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAllUsers(page, limit int) (*models.PaginatedResponse[models.User], error) {
	args := m.Called(page, limit)
	result, _ := args.Get(0).(*models.PaginatedResponse[models.User])
	return result, args.Error(1)
}

func (m *MockUserService) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	result, _ := args.Get(0).(*models.User)
	return result, args.Error(1)
}

func (m *MockUserService) UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(id, req)
	result, _ := args.Get(0).(*models.User)
	return result, args.Error(1)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// ────────────────────────────────────────────────────────────────
// MockPostService  (implements services.IPostService)
// ────────────────────────────────────────────────────────────────

type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) CreatePost(post models.Post, userID uint) (*models.Post, error) {
	args := m.Called(post, userID)
	result, _ := args.Get(0).(*models.Post)
	return result, args.Error(1)
}

func (m *MockPostService) GetAllPosts(page, limit int) (*models.PaginatedResponse[models.Post], error) {
	args := m.Called(page, limit)
	result, _ := args.Get(0).(*models.PaginatedResponse[models.Post])
	return result, args.Error(1)
}

func (m *MockPostService) GetPostByID(id uint) (*models.Post, error) {
	args := m.Called(id)
	result, _ := args.Get(0).(*models.Post)
	return result, args.Error(1)
}

func (m *MockPostService) GetUserPosts(userID uint, page, limit int) (*models.PaginatedResponse[models.Post], error) {
	args := m.Called(userID, page, limit)
	result, _ := args.Get(0).(*models.PaginatedResponse[models.Post])
	return result, args.Error(1)
}

func (m *MockPostService) UpdatePost(id uint, post models.Post, userID uint) (*models.Post, error) {
	args := m.Called(id, post, userID)
	result, _ := args.Get(0).(*models.Post)
	return result, args.Error(1)
}

func (m *MockPostService) DeletePost(id uint, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}
