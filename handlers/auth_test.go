package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"awesomeProject/errors"
	"awesomeProject/handlers"
	"awesomeProject/mocks"
	"awesomeProject/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthRouter(svc *mocks.MockAuthService) *gin.Engine {
	r := gin.New()
	h := handlers.NewAuthHandler(svc)
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.Refresh)
	r.POST("/auth/logout", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	}, h.Logout)
	return r
}

// ────────────────────────────────────────────────────────────────
// Register
// ────────────────────────────────────────────────────────────────

func TestRegisterHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockAuthService{}
	reqBody := models.RegisterRequest{Name: "João Silva", Email: "joao@example.com", Password: "senha123"}
	respBody := &models.AuthResponse{ID: 1, Name: "João Silva", Email: "joao@example.com", Token: "tok", RefreshToken: "ref"}

	mockSvc.On("Register", reqBody).Return(respBody, nil)

	r := setupAuthRouter(mockSvc)
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var got models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "João Silva", got.Name)
	assert.Equal(t, "tok", got.Token)

	mockSvc.AssertExpectations(t)
}

func TestRegisterHandler_InvalidBody(t *testing.T) {
	mockSvc := &mocks.MockAuthService{}
	r := setupAuthRouter(mockSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{bad json}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterHandler_EmailConflict(t *testing.T) {
	mockSvc := &mocks.MockAuthService{}
	reqBody := models.RegisterRequest{Name: "João Silva", Email: "joao@example.com", Password: "senha123"}

	mockSvc.On("Register", reqBody).Return(nil, errors.ErrConflict("email já cadastrado"))

	r := setupAuthRouter(mockSvc)
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockSvc.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// Login
// ────────────────────────────────────────────────────────────────

func TestLoginHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockAuthService{}
	reqBody := models.LoginRequest{Email: "joao@example.com", Password: "senha123"}
	respBody := &models.AuthResponse{ID: 1, Name: "João", Email: "joao@example.com", Token: "tok", RefreshToken: "ref"}

	mockSvc.On("Login", reqBody).Return(respBody, nil)

	r := setupAuthRouter(mockSvc)
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "tok", got.Token)

	mockSvc.AssertExpectations(t)
}

func TestLoginHandler_Unauthorized(t *testing.T) {
	mockSvc := &mocks.MockAuthService{}
	reqBody := models.LoginRequest{Email: "joao@example.com", Password: "errada"}

	mockSvc.On("Login", reqBody).Return(nil, errors.ErrUnauthorized("email ou senha inválidos"))

	r := setupAuthRouter(mockSvc)
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// Refresh
// ────────────────────────────────────────────────────────────────

func TestRefreshHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockAuthService{}
	respBody := &models.AuthResponse{Token: "new-tok", RefreshToken: "new-ref"}

	mockSvc.On("RefreshTokens", "old-refresh").Return(respBody, nil)

	r := setupAuthRouter(mockSvc)
	reqBody := models.RefreshRequest{RefreshToken: "old-refresh"}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "new-tok", got.Token)

	mockSvc.AssertExpectations(t)
}

func TestRefreshHandler_InvalidToken(t *testing.T) {
	mockSvc := &mocks.MockAuthService{}

	mockSvc.On("RefreshTokens", "bad-token").Return(nil, errors.ErrUnauthorized("refresh token inválido"))

	r := setupAuthRouter(mockSvc)
	reqBody := models.RefreshRequest{RefreshToken: "bad-token"}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// Logout
// ────────────────────────────────────────────────────────────────

func TestLogoutHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockAuthService{}
	mockSvc.On("Logout", uint(1)).Return(nil)

	r := setupAuthRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}
