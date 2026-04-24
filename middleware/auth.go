package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"awesomeProject/errors"
	"awesomeProject/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TokenValidator is implemented by any service capable of validating JWT tokens.
type TokenValidator interface {
	ValidateToken(tokenString string) (*services.Claims, error)
}

// RequireAuthWithValidator creates an auth middleware using the provided TokenValidator.
// Prefer this over RequireAuth for better testability.
func RequireAuthWithValidator(v TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errors.HandleError(c, errors.ErrUnauthorized("token não fornecido"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errors.HandleError(c, errors.ErrUnauthorized("formato de token inválido"))
			c.Abort()
			return
		}

		claims, err := v.ValidateToken(parts[1])
		if err != nil {
			errors.HandleError(c, err)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// RequireAuth creates auth middleware with a default AuthService (nil repos are safe
// because ValidateToken only reads config.AppConfig.JWTSecret).
func RequireAuth() gin.HandlerFunc {
	return RequireAuthWithValidator(&services.AuthService{})
}

// StructuredLogging middleware provides structured logging with request ID, timestamp, and level
func StructuredLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a unique request ID
		requestID := uuid.New().String()
		c.Set("requestID", requestID)

		// Start time for request duration calculation
		startTime := time.Now()

		// Log request details
		slog.Info("incoming request",
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Time("timestamp", startTime),
		)

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(startTime)

		// Log response details
		statusCode := c.Writer.Status()
		var level slog.Level

		// Set log level based on HTTP status code
		switch {
		case statusCode >= 500:
			level = slog.LevelError
		case statusCode >= 400:
			level = slog.LevelWarn
		default:
			level = slog.LevelInfo
		}

		slog.Log(c.Request.Context(), level, "request completed",
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status_code", statusCode),
			slog.String("status_text", http.StatusText(statusCode)),
			slog.Duration("duration", duration),
			slog.Int("response_size", c.Writer.Size()),
			slog.Time("timestamp", time.Now()),
		)
	}
}
