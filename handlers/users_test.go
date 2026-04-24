package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"awesomeProject/errors"
	"awesomeProject/handlers"
	"awesomeProject/mocks"
	"awesomeProject/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupUserRouter(svc *mocks.MockUserService) *gin.Engine {
	r := gin.New()
	h := handlers.NewUserHandler(svc)
	r.GET("/users/", h.GetAllUsers)
	r.GET("/users/:id", h.GetUserByID)

	// Routes that need auth context
	r.PUT("/users/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	}, h.UpdateUser)
	r.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	}, h.DeleteUser)
	return r
}

// ────────────────────────────────────────────────────────────────
// GetAllUsers
// ────────────────────────────────────────────────────────────────

func TestGetAllUsersHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockUserService{}
	resp := &models.PaginatedResponse[models.User]{
		Data:       []models.User{{ID: 1, Name: "Alice"}},
		Page:       1,
		Limit:      10,
		Total:      1,
		TotalPages: 1,
	}
	mockSvc.On("GetAllUsers", 1, 10).Return(resp, nil)

	r := setupUserRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/?page=1&limit=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got models.PaginatedResponse[models.User]
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got.Data, 1)

	mockSvc.AssertExpectations(t)
}

func TestGetAllUsersHandler_ServiceError(t *testing.T) {
	mockSvc := &mocks.MockUserService{}
	mockSvc.On("GetAllUsers", 1, 10).Return(nil, errors.ErrInternal("db error"))

	r := setupUserRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/?page=1&limit=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// GetUserByID
// ────────────────────────────────────────────────────────────────

func TestGetUserByIDHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockUserService{}
	user := &models.User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	mockSvc.On("GetUserByID", uint(1)).Return(user, nil)

	r := setupUserRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got models.User
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "Alice", got.Name)

	mockSvc.AssertExpectations(t)
}

func TestGetUserByIDHandler_NotFound(t *testing.T) {
	mockSvc := &mocks.MockUserService{}
	mockSvc.On("GetUserByID", uint(99)).Return(nil, errors.ErrNotFound("usuário não encontrado"))

	r := setupUserRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetUserByIDHandler_InvalidID(t *testing.T) {
	mockSvc := &mocks.MockUserService{}
	r := setupUserRouter(mockSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ────────────────────────────────────────────────────────────────
// UpdateUser
// ────────────────────────────────────────────────────────────────

func TestUpdateUserHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockUserService{}
	newName := "Alice Updated"
	reqBody := models.UpdateUserRequest{Name: &newName}
	updatedUser := &models.User{ID: 1, Name: "Alice Updated", Email: "alice@example.com"}

	mockSvc.On("UpdateUser", uint(1), reqBody).Return(updatedUser, nil)

	r := setupUserRouter(mockSvc)
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got models.User
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "Alice Updated", got.Name)

	mockSvc.AssertExpectations(t)
}

func TestUpdateUserHandler_Forbidden_WrongOwner(t *testing.T) {
	// The auth middleware sets userID=1, but the route is /users/2
	mockSvc := &mocks.MockUserService{}
	r := setupUserRouter(mockSvc) // userID always = 1

	body, _ := json.Marshal(map[string]string{"name": "Hacker"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/users/2", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ────────────────────────────────────────────────────────────────
// DeleteUser
// ────────────────────────────────────────────────────────────────

func TestDeleteUserHandler_Success(t *testing.T) {
	mockSvc := &mocks.MockUserService{}
	mockSvc.On("DeleteUser", uint(1)).Return(nil)

	r := setupUserRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteUserHandler_Forbidden_WrongOwner(t *testing.T) {
	mockSvc := &mocks.MockUserService{}
	r := setupUserRouter(mockSvc) // userID always = 1

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users/2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteUserHandler_NotFound(t *testing.T) {
	mockSvc := &mocks.MockUserService{}
	mockSvc.On("DeleteUser", uint(1)).Return(errors.ErrNotFound("usuário não encontrado"))

	r := setupUserRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

// compile-time interface checks
var _ = fmt.Sprintf // keep import
