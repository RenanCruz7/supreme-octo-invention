package handlers

import (
	"net/http"

	"awesomeProject/errors"
	"awesomeProject/models"
	"awesomeProject/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc services.IAuthService
}

func NewAuthHandler(svc services.IAuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrBadRequestWithErr("erro ao processar registro", err))
		return
	}

	resp, err := h.svc.Register(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrBadRequestWithErr("erro ao processar login", err))
		return
	}

	resp, err := h.svc.Login(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrBadRequestWithErr("refresh_token é obrigatório", err))
		return
	}

	resp, err := h.svc.RefreshTokens(req.RefreshToken)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	if err := h.svc.Logout(userID); err != nil {
		errors.HandleError(c, errors.ErrInternalWithErr("erro ao fazer logout", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout realizado com sucesso"})
}
