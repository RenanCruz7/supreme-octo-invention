package handlers

import (
	"net/http"
	"strconv"

	"awesomeProject/errors"
	"awesomeProject/models"
	"awesomeProject/services"

	"github.com/gin-gonic/gin"
)

var userService = &services.UserService{}

func GetAllUsers(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	users, err := userService.GetAllUsers(page, limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest("ID inválido"))
		return
	}

	user, err := userService.GetUserByID(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest("ID inválido"))
		return
	}

	userID, exists := c.Get("userID")
	if !exists || userID.(uint) != uint(id) {
		errors.HandleError(c, errors.ErrForbidden("você pode apenas atualizar sua própria conta"))
		return
	}

	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrBadRequestWithErr("erro ao processar usuário", err))
		return
	}

	user, err := userService.UpdateUser(uint(id), req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest("ID inválido"))
		return
	}

	userID, exists := c.Get("userID")
	if !exists || userID.(uint) != uint(id) {
		errors.HandleError(c, errors.ErrForbidden("você pode apenas deletar sua própria conta"))
		return
	}

	err = userService.DeleteUser(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
