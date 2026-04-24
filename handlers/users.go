package handlers

import (
	"net/http"
	"strconv"

	"awesomeProject/errors"
	"awesomeProject/models"
	"awesomeProject/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc services.IUserService
}

func NewUserHandler(svc services.IUserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, limit := parsePagination(c)

	result, err := h.svc.GetAllUsers(page, limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest("ID inválido"))
		return
	}

	user, err := h.svc.GetUserByID(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
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

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrBadRequestWithErr("erro ao processar usuário", err))
		return
	}

	user, err := h.svc.UpdateUser(uint(id), req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
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

	err = h.svc.DeleteUser(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
