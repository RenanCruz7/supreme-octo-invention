package services_test

import (
	"fmt"
	"testing"

	"awesomeProject/mocks"
	"awesomeProject/models"
	"awesomeProject/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ────────────────────────────────────────────────────────────────
// GetAllUsers
// ────────────────────────────────────────────────────────────────

func TestGetAllUsers_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	users := []models.User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}
	userRepo.On("GetAllUsers", 1, 10).Return(users, nil)
	userRepo.On("CountAllUsers").Return(2, nil)

	resp, err := svc.GetAllUsers(1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, 2, resp.Total)
	assert.Equal(t, 1, resp.TotalPages)

	userRepo.AssertExpectations(t)
}

func TestGetAllUsers_DefaultPagination(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	// page=0 and limit=0 should be corrected to page=1, limit=10
	userRepo.On("GetAllUsers", 1, 10).Return([]models.User{}, nil)
	userRepo.On("CountAllUsers").Return(0, nil)

	resp, err := svc.GetAllUsers(0, 0)

	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.Limit)

	userRepo.AssertExpectations(t)
}

func TestGetAllUsers_RepoError(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	userRepo.On("GetAllUsers", 1, 10).Return(nil, fmt.Errorf("db error"))

	resp, err := svc.GetAllUsers(1, 10)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "INTERNAL_ERROR")

	userRepo.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// GetUserByID
// ────────────────────────────────────────────────────────────────

func TestGetUserByID_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	user := &models.User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	userRepo.On("GetUserByID", uint(1)).Return(user, nil)

	result, err := svc.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)

	userRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	userRepo.On("GetUserByID", uint(99)).Return(nil, fmt.Errorf("not found"))

	result, err := svc.GetUserByID(99)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NOT_FOUND")

	userRepo.AssertExpectations(t)
}

func TestGetUserByID_InvalidID(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	result, err := svc.GetUserByID(0)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BAD_REQUEST")
}

// ────────────────────────────────────────────────────────────────
// UpdateUser
// ────────────────────────────────────────────────────────────────

func TestUpdateUser_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	existing := &models.User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	userRepo.On("GetUserByID", uint(1)).Return(existing, nil)
	userRepo.On("UpdateUser", mock.AnythingOfType("User")).Return(nil)

	newName := "Alice Updated"
	req := models.UpdateUserRequest{Name: &newName}

	result, err := svc.UpdateUser(1, req)

	assert.NoError(t, err)
	assert.Equal(t, "Alice Updated", result.Name)

	userRepo.AssertExpectations(t)
}

func TestUpdateUser_EmailConflict(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	existing := &models.User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	otherUser := &models.User{ID: 2, Email: "taken@example.com"}
	userRepo.On("GetUserByID", uint(1)).Return(existing, nil)
	userRepo.On("GetUserByEmail", "taken@example.com").Return(otherUser, nil)

	takenEmail := "taken@example.com"
	req := models.UpdateUserRequest{Email: &takenEmail}

	result, err := svc.UpdateUser(1, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CONFLICT")

	userRepo.AssertExpectations(t)
}

func TestUpdateUser_NotFound(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	userRepo.On("GetUserByID", uint(1)).Return(nil, fmt.Errorf("not found"))

	req := models.UpdateUserRequest{}
	result, err := svc.UpdateUser(1, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NOT_FOUND")

	userRepo.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// DeleteUser
// ────────────────────────────────────────────────────────────────

func TestDeleteUser_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	user := &models.User{ID: 1}
	userRepo.On("GetUserByID", uint(1)).Return(user, nil)
	userRepo.On("DeleteUser", uint(1)).Return(nil)

	err := svc.DeleteUser(1)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	userRepo.On("GetUserByID", uint(99)).Return(nil, fmt.Errorf("not found"))

	err := svc.DeleteUser(99)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NOT_FOUND")

	userRepo.AssertExpectations(t)
}

func TestDeleteUser_InvalidID(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	svc := services.NewUserService(userRepo)

	err := svc.DeleteUser(0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BAD_REQUEST")
}
