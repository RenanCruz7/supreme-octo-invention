package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"awesomeProject/errors"
	"awesomeProject/handlers"
	"awesomeProject/mocks"
	"awesomeProject/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupPostRouter(svc *mocks.MockPostService) *gin.Engine {
	r := gin.New()
	h := handlers.NewPostHandler(svc)
	r.GET("/posts/", h.GetAllPosts)
	r.GET("/posts/:id", h.GetPostByID)
	r.GET("/users/:id/posts", h.GetUserPosts)

	// Protected routes – pre-set userID=1 in context
	authed := func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	}
	r.POST("/posts/", authed, h.CreatePost)
	r.PUT("/posts/:id", authed, h.UpdatePost)
	r.DELETE("/posts/:id", authed, h.DeletePost)
	return r
}

// ────────────────────────────────────────────────────────────────
// GetAllPosts
// ────────────────────────────────────────────────────────────────

func TestGetAllPostsHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	resp := &models.PaginatedResponse[models.Post]{
		Data:       []models.Post{{ID: 1, Title: "Post 1"}, {ID: 2, Title: "Post 2"}},
		Page:       1,
		Limit:      10,
		Total:      2,
		TotalPages: 1,
	}
	mockSvc.On("GetAllPosts", 1, 10).Return(resp, nil)

	r := setupPostRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts/?page=1&limit=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got models.PaginatedResponse[models.Post]
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got.Data, 2)

	mockSvc.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// GetPostByID
// ────────────────────────────────────────────────────────────────

func TestGetPostByIDHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	post := &models.Post{ID: 5, Title: "Meu Post", UserID: 1}
	mockSvc.On("GetPostByID", uint(5)).Return(post, nil)

	r := setupPostRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts/5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got models.Post
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, uint(5), got.ID)

	mockSvc.AssertExpectations(t)
}

func TestGetPostByIDHandler_NotFound(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	mockSvc.On("GetPostByID", uint(99)).Return(nil, errors.ErrNotFound("post não encontrado"))

	r := setupPostRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// GetUserPosts
// ────────────────────────────────────────────────────────────────

func TestGetUserPostsHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	resp := &models.PaginatedResponse[models.Post]{
		Data: []models.Post{{ID: 1, UserID: 1}},
	}
	mockSvc.On("GetUserPosts", uint(1), 1, 10).Return(resp, nil)

	r := setupPostRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/1/posts", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetUserPostsHandler_UserNotFound(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	mockSvc.On("GetUserPosts", uint(99), 1, 10).Return(nil, errors.ErrNotFound("usuário não encontrado"))

	r := setupPostRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/99/posts", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// CreatePost
// ────────────────────────────────────────────────────────────────

func TestCreatePostHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	postInput := models.Post{Title: "Novo Post", Body: "Conteúdo do meu novo post aqui"}
	createdPost := &models.Post{ID: 1, Title: "Novo Post", Body: "Conteúdo do meu novo post aqui", UserID: 1}

	mockSvc.On("CreatePost", postInput, uint(1)).Return(createdPost, nil)

	r := setupPostRouter(mockSvc)
	body, _ := json.Marshal(postInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/posts/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var got models.Post
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, uint(1), got.ID)

	mockSvc.AssertExpectations(t)
}

func TestCreatePostHandler_InvalidBody(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	r := setupPostRouter(mockSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/posts/", bytes.NewBufferString(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ────────────────────────────────────────────────────────────────
// UpdatePost
// ────────────────────────────────────────────────────────────────

func TestUpdatePostHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	postInput := models.Post{Title: "Título Atualizado", Body: "Corpo atualizado com novo conteúdo"}
	updatedPost := &models.Post{ID: 1, Title: "Título Atualizado", Body: "Corpo atualizado com novo conteúdo", UserID: 1}

	mockSvc.On("UpdatePost", uint(1), postInput, uint(1)).Return(updatedPost, nil)

	r := setupPostRouter(mockSvc)
	body, _ := json.Marshal(postInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/posts/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got models.Post
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "Título Atualizado", got.Title)

	mockSvc.AssertExpectations(t)
}

func TestUpdatePostHandler_Forbidden(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	postInput := models.Post{Title: "Hack", Body: "Tentativa de hack no sistema"}

	mockSvc.On("UpdatePost", uint(2), postInput, uint(1)).Return(nil, errors.ErrForbidden("você não tem permissão"))

	r := setupPostRouter(mockSvc)
	body, _ := json.Marshal(postInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/posts/2", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockSvc.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// DeletePost
// ────────────────────────────────────────────────────────────────

func TestDeletePostHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	mockSvc.On("DeletePost", uint(1), uint(1)).Return(nil)

	r := setupPostRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/posts/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeletePostHandler_Forbidden(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	mockSvc.On("DeletePost", uint(2), uint(1)).Return(errors.ErrForbidden("você não tem permissão para deletar este post"))

	r := setupPostRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/posts/2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeletePostHandler_NotFound(t *testing.T) {
	mockSvc := &mocks.MockPostService{}
	mockSvc.On("DeletePost", uint(99), uint(1)).Return(errors.ErrNotFound("post não encontrado"))

	r := setupPostRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/posts/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

// keep mock import used
var _ = mock.Anything
