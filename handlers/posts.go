package handlers

import (
	"net/http"
	"strconv"

	"awesomeProject/errors"
	"awesomeProject/models"
	"awesomeProject/services"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	svc services.IPostService
}

func NewPostHandler(svc services.IPostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		errors.HandleError(c, errors.ErrBadRequestWithErr("erro ao processar post", err))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		errors.HandleError(c, errors.ErrUnauthorized("usuário não autenticado"))
		return
	}

	createdPost, err := h.svc.CreatePost(post, userID.(uint))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdPost)
}

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	page, limit := parsePagination(c)

	result, err := h.svc.GetAllPosts(page, limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest("ID inválido"))
		return
	}

	post, err := h.svc.GetPostByID(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest("User ID inválido"))
		return
	}

	page, limit := parsePagination(c)

	result, err := h.svc.GetUserPosts(uint(userID), page, limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest("ID inválido"))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		errors.HandleError(c, errors.ErrUnauthorized("usuário não autenticado"))
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		errors.HandleError(c, errors.ErrBadRequestWithErr("erro ao processar post", err))
		return
	}

	updatedPost, err := h.svc.UpdatePost(uint(id), post, userID.(uint))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest("ID inválido"))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		errors.HandleError(c, errors.ErrUnauthorized("usuário não autenticado"))
		return
	}

	err = h.svc.DeletePost(uint(id), userID.(uint))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
