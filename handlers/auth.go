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
