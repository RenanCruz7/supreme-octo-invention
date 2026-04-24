package services_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"awesomeProject/config"
	"awesomeProject/mocks"
	"awesomeProject/models"
	"awesomeProject/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	config.AppConfig = &config.Config{
		JWTSecret: "test_secret_key_for_testing_123456",
	}
	os.Exit(m.Run())
}

// ────────────────────────────────────────────────────────────────
// Register
// ────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenRepo := &mocks.MockRefreshTokenRepository{}
	svc := services.NewAuthService(userRepo, tokenRepo)

	req := models.RegisterRequest{Name: "João Silva", Email: "joao@example.com", Password: "senha123"}

	userRepo.On("GetUserByEmail", "joao@example.com").Return(nil, fmt.Errorf("not found"))
	userRepo.On("CreateUser", mock.AnythingOfType("User")).Return(uint(1), nil)
	tokenRepo.On("CreateRefreshToken", mock.AnythingOfType("RefreshToken")).Return(nil)

	resp, err := svc.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "João Silva", resp.Name)
	assert.Equal(t, "joao@example.com", resp.Email)
	assert.NotEmpty(t, resp.Token)
	assert.NotEmpty(t, resp.RefreshToken)

	userRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenRepo := &mocks.MockRefreshTokenRepository{}
	svc := services.NewAuthService(userRepo, tokenRepo)

	existing := &models.User{ID: 1, Email: "joao@example.com"}
	userRepo.On("GetUserByEmail", "joao@example.com").Return(existing, nil)

	req := models.RegisterRequest{Name: "João Silva", Email: "joao@example.com", Password: "senha123"}
	resp, err := svc.Register(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CONFLICT")

	userRepo.AssertExpectations(t)
}

func TestRegister_CreateUserError(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenRepo := &mocks.MockRefreshTokenRepository{}
	svc := services.NewAuthService(userRepo, tokenRepo)

	userRepo.On("GetUserByEmail", "joao@example.com").Return(nil, fmt.Errorf("not found"))
	userRepo.On("CreateUser", mock.AnythingOfType("User")).Return(uint(0), fmt.Errorf("db error"))

	req := models.RegisterRequest{Name: "João Silva", Email: "joao@example.com", Password: "senha123"}
	resp, err := svc.Register(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "INTERNAL_ERROR")

	userRepo.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// Login
// ────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenRepo := &mocks.MockRefreshTokenRepository{}
	svc := services.NewAuthService(userRepo, tokenRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.MinCost)
	user := &models.User{ID: 1, Name: "João", Email: "joao@example.com", Password: string(hashed)}

	userRepo.On("GetUserByEmail", "joao@example.com").Return(user, nil)
	tokenRepo.On("CreateRefreshToken", mock.AnythingOfType("RefreshToken")).Return(nil)

	req := models.LoginRequest{Email: "joao@example.com", Password: "senha123"}
	resp, err := svc.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	assert.NotEmpty(t, resp.Token)
	assert.NotEmpty(t, resp.RefreshToken)

	userRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenRepo := &mocks.MockRefreshTokenRepository{}
	svc := services.NewAuthService(userRepo, tokenRepo)

	userRepo.On("GetUserByEmail", "inexistente@example.com").Return(nil, fmt.Errorf("not found"))

	req := models.LoginRequest{Email: "inexistente@example.com", Password: "senha123"}
	resp, err := svc.Login(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNAUTHORIZED")

	userRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenRepo := &mocks.MockRefreshTokenRepository{}
	svc := services.NewAuthService(userRepo, tokenRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("senhaCorreta"), bcrypt.MinCost)
	user := &models.User{ID: 1, Email: "joao@example.com", Password: string(hashed)}

	userRepo.On("GetUserByEmail", "joao@example.com").Return(user, nil)

	req := models.LoginRequest{Email: "joao@example.com", Password: "senhaErrada"}
	resp, err := svc.Login(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNAUTHORIZED")

	userRepo.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// Logout
// ────────────────────────────────────────────────────────────────

func TestLogout_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenRepo := &mocks.MockRefreshTokenRepository{}
	svc := services.NewAuthService(userRepo, tokenRepo)

	tokenRepo.On("RevokeAllUserTokens", uint(1)).Return(nil)

	err := svc.Logout(1)

	assert.NoError(t, err)
	tokenRepo.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// RefreshTokens
// ────────────────────────────────────────────────────────────────

func TestRefreshTokens_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenRepo := &mocks.MockRefreshTokenRepository{}
	svc := services.NewAuthService(userRepo, tokenRepo)

	rt := &models.RefreshToken{
		UserID:    1,
		Token:     "valid-refresh-token",
		ExpiresAt: time.Now().Add(time.Hour),
		Revoked:   false,
	}
	user := &models.User{ID: 1, Name: "João", Email: "joao@example.com"}

	tokenRepo.On("GetRefreshToken", "valid-refresh-token").Return(rt, nil)
	tokenRepo.On("RevokeRefreshToken", "valid-refresh-token").Return(nil)
	userRepo.On("GetUserByID", uint(1)).Return(user, nil)
	tokenRepo.On("CreateRefreshToken", mock.AnythingOfType("RefreshToken")).Return(nil)

	resp, err := svc.RefreshTokens("valid-refresh-token")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.NotEmpty(t, resp.RefreshToken)

	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshTokens_InvalidToken(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenRepo := &mocks.MockRefreshTokenRepository{}
	svc := services.NewAuthService(userRepo, tokenRepo)

	tokenRepo.On("GetRefreshToken", "bad-token").Return(nil, fmt.Errorf("not found"))

	resp, err := svc.RefreshTokens("bad-token")

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNAUTHORIZED")

	tokenRepo.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// GenerateToken / ValidateToken
// ────────────────────────────────────────────────────────────────

func TestGenerateAndValidateToken(t *testing.T) {
	svc := &services.AuthService{}

	token, err := svc.GenerateToken(42, "test@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := svc.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	svc := &services.AuthService{}
	_, err := svc.ValidateToken("invalid.token.here")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNAUTHORIZED")
}
