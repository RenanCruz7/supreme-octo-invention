package handlers

import (
	"net/http"

	"awesomeProject/errors"
	"awesomeProject/models"
	"awesomeProject/services"

	"github.com/gin-gonic/gin"
)

var authService = &services.AuthService{}

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrBadRequestWithErr("erro ao processar registro", err))
		return
	}

	resp, err := authService.Register(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrBadRequestWithErr("erro ao processar login", err))
		return
	}

	resp, err := authService.Login(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Refresh recebe um refresh token e retorna um novo par access+refresh token
func Refresh(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrBadRequestWithErr("refresh_token é obrigatório", err))
		return
	}

	resp, err := authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Logout revoga todos os refresh tokens do usuário autenticado
func Logout(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	if err := authService.Logout(userID); err != nil {
		errors.HandleError(c, errors.ErrInternalWithErr("erro ao fazer logout", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout realizado com sucesso"})
}
