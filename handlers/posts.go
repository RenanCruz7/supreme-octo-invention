package handlers

import (
	"net/http"
	"strconv"

	"awesomeProject/errors"
	"awesomeProject/models"
	"awesomeProject/services"

	"github.com/gin-gonic/gin"
)

var postService = &services.PostService{}

func CreatePost(c *gin.Context) {
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

	createdPost, err := postService.CreatePost(post, userID.(uint))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdPost)
}

func GetAllPosts(c *gin.Context) {
	posts, err := postService.GetAllPosts()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func GetPostByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest("ID inválido"))
		return
	}

	post, err := postService.GetPostByID(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

func GetUserPosts(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest("User ID inválido"))
		return
	}

	posts, err := postService.GetUserPosts(uint(userID))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func UpdatePost(c *gin.Context) {
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

	updatedPost, err := postService.UpdatePost(uint(id), post, userID.(uint))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

func DeletePost(c *gin.Context) {
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

	err = postService.DeletePost(uint(id), userID.(uint))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
