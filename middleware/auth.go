package middleware

import (
	"strings"

	"awesomeProject/errors"
	"awesomeProject/services"

	"github.com/gin-gonic/gin"
)

var authService = &services.AuthService{}

func RequireAuth() gin.HandlerFunc {
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

		tokenString := parts[1]

		claims, err := authService.ValidateToken(tokenString)
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
