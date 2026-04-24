package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"awesomeProject/config"
	"awesomeProject/middleware"
	"awesomeProject/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
	config.AppConfig = &config.Config{
		JWTSecret: "test_secret_key_for_testing_123456",
	}
}

func newValidToken(userID uint, email string) string {
	svc := &services.AuthService{}
	tok, _ := svc.GenerateToken(userID, email)
	return tok
}

func setupAuthRouter(v middleware.TokenValidator) *gin.Engine {
	r := gin.New()
	r.GET("/protected", middleware.RequireAuthWithValidator(v), func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})
	return r
}

func TestRequireAuth_NoHeader(t *testing.T) {
	svc := &services.AuthService{}
	r := setupAuthRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireAuth_InvalidHeaderFormat(t *testing.T) {
	svc := &services.AuthService{}
	r := setupAuthRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireAuth_ValidToken_SetsUserID(t *testing.T) {
	svc := &services.AuthService{}
	token := newValidToken(42, "test@example.com")
	r := setupAuthRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "42")
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	svc := &services.AuthService{}
	r := setupAuthRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer this.is.not.valid")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireAuth_WrongSecret(t *testing.T) {
	// Token generated with a different secret
	other := &services.AuthService{}
	originalSecret := config.AppConfig.JWTSecret
	config.AppConfig.JWTSecret = "other_secret"
	badToken, _ := other.GenerateToken(1, "x@x.com")
	config.AppConfig.JWTSecret = originalSecret // restore

	svc := &services.AuthService{}
	r := setupAuthRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+badToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
