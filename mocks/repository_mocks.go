package mocks

import (
	"awesomeProject/models"

	"github.com/stretchr/testify/mock"
)

// ────────────────────────────────────────────────────────────────
// MockUserRepository
// ────────────────────────────────────────────────────────────────

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user models.User) (uint, error) {
	args := m.Called(user)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers(page, limit int) ([]models.User, error) {
	args := m.Called(page, limit)
	result, _ := args.Get(0).([]models.User)
	return result, args.Error(1)
}

func (m *MockUserRepository) CountAllUsers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	result, _ := args.Get(0).(*models.User)
	return result, args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	result, _ := args.Get(0).(*models.User)
	return result, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// ────────────────────────────────────────────────────────────────
// MockPostRepository
// ────────────────────────────────────────────────────────────────

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) CreatePost(post models.Post) (uint, error) {
	args := m.Called(post)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockPostRepository) GetAllPosts(page, limit int) ([]models.Post, error) {
	args := m.Called(page, limit)
	result, _ := args.Get(0).([]models.Post)
	return result, args.Error(1)
}

func (m *MockPostRepository) CountAllPosts() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockPostRepository) GetPostByID(id uint) (*models.Post, error) {
	args := m.Called(id)
	result, _ := args.Get(0).(*models.Post)
	return result, args.Error(1)
}

func (m *MockPostRepository) GetPostsByUserID(userID uint, page, limit int) ([]models.Post, error) {
	args := m.Called(userID, page, limit)
	result, _ := args.Get(0).([]models.Post)
	return result, args.Error(1)
}

func (m *MockPostRepository) CountPostsByUserID(userID uint) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockPostRepository) UpdatePost(post models.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) DeletePost(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// ────────────────────────────────────────────────────────────────
// MockRefreshTokenRepository
// ────────────────────────────────────────────────────────────────

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) CreateRefreshToken(token models.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	args := m.Called(token)
	result, _ := args.Get(0).(*models.RefreshToken)
	return result, args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeRefreshToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllUserTokens(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}
