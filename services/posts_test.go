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
// CreatePost
// ────────────────────────────────────────────────────────────────

func TestCreatePost_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	user := &models.User{ID: 1}
	post := models.Post{Title: "Meu Post", Body: "Corpo do post com conteúdo"}

	userRepo.On("GetUserByID", uint(1)).Return(user, nil)
	postRepo.On("CreatePost", mock.AnythingOfType("Post")).Return(uint(10), nil)

	result, err := svc.CreatePost(post, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(10), result.ID)
	assert.Equal(t, uint(1), result.UserID)

	userRepo.AssertExpectations(t)
	postRepo.AssertExpectations(t)
}

func TestCreatePost_UserNotFound(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	userRepo.On("GetUserByID", uint(99)).Return(nil, fmt.Errorf("not found"))

	post := models.Post{Title: "Post", Body: "Corpo do post"}
	result, err := svc.CreatePost(post, 99)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NOT_FOUND")

	userRepo.AssertExpectations(t)
}

func TestCreatePost_UnauthenticatedUser(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	post := models.Post{Title: "Post", Body: "Corpo"}
	result, err := svc.CreatePost(post, 0)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNAUTHORIZED")
}

// ────────────────────────────────────────────────────────────────
// GetAllPosts
// ────────────────────────────────────────────────────────────────

func TestGetAllPosts_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	posts := []models.Post{
		{ID: 1, Title: "Post 1", UserID: 1},
		{ID: 2, Title: "Post 2", UserID: 1},
	}
	postRepo.On("GetAllPosts", 1, 10).Return(posts, nil)
	postRepo.On("CountAllPosts").Return(2, nil)

	resp, err := svc.GetAllPosts(1, 10)

	assert.NoError(t, err)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, 2, resp.Total)

	postRepo.AssertExpectations(t)
}

func TestGetAllPosts_LimitClamped(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	// limit=200 should be clamped to 10
	postRepo.On("GetAllPosts", 1, 10).Return([]models.Post{}, nil)
	postRepo.On("CountAllPosts").Return(0, nil)

	resp, err := svc.GetAllPosts(1, 200)

	assert.NoError(t, err)
	assert.Equal(t, 10, resp.Limit)

	postRepo.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// GetPostByID
// ────────────────────────────────────────────────────────────────

func TestGetPostByID_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	post := &models.Post{ID: 5, Title: "Post", UserID: 1}
	postRepo.On("GetPostByID", uint(5)).Return(post, nil)

	result, err := svc.GetPostByID(5)

	assert.NoError(t, err)
	assert.Equal(t, uint(5), result.ID)

	postRepo.AssertExpectations(t)
}

func TestGetPostByID_NotFound(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	postRepo.On("GetPostByID", uint(99)).Return(nil, fmt.Errorf("not found"))

	result, err := svc.GetPostByID(99)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NOT_FOUND")

	postRepo.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// UpdatePost
// ────────────────────────────────────────────────────────────────

func TestUpdatePost_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	existing := &models.Post{ID: 1, Title: "Original", Body: "Corpo original", UserID: 1}
	postRepo.On("GetPostByID", uint(1)).Return(existing, nil)
	postRepo.On("UpdatePost", mock.AnythingOfType("Post")).Return(nil)

	update := models.Post{Title: "Novo Título", Body: "Novo corpo do post"}
	result, err := svc.UpdatePost(1, update, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Novo Título", result.Title)
	assert.Equal(t, "Novo corpo do post", result.Body)

	postRepo.AssertExpectations(t)
}

func TestUpdatePost_Forbidden(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	existing := &models.Post{ID: 1, Title: "Post", Body: "Corpo", UserID: 2}
	postRepo.On("GetPostByID", uint(1)).Return(existing, nil)

	update := models.Post{Title: "Hack", Body: "Tentativa de hack"}
	result, err := svc.UpdatePost(1, update, 1) // userID=1 trying to edit ownerID=2 post

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FORBIDDEN")

	postRepo.AssertExpectations(t)
}

func TestUpdatePost_NotFound(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	postRepo.On("GetPostByID", uint(99)).Return(nil, fmt.Errorf("not found"))

	result, err := svc.UpdatePost(99, models.Post{}, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NOT_FOUND")

	postRepo.AssertExpectations(t)
}

// ────────────────────────────────────────────────────────────────
// DeletePost
// ────────────────────────────────────────────────────────────────

func TestDeletePost_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	post := &models.Post{ID: 1, UserID: 1}
	postRepo.On("GetPostByID", uint(1)).Return(post, nil)
	postRepo.On("DeletePost", uint(1)).Return(nil)

	err := svc.DeletePost(1, 1)

	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestDeletePost_Forbidden(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	post := &models.Post{ID: 1, UserID: 2}
	postRepo.On("GetPostByID", uint(1)).Return(post, nil)

	err := svc.DeletePost(1, 1) // userID=1 trying to delete ownerID=2 post

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FORBIDDEN")

	postRepo.AssertExpectations(t)
}

func TestDeletePost_InvalidID(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	err := svc.DeletePost(0, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BAD_REQUEST")
}

// ────────────────────────────────────────────────────────────────
// GetUserPosts
// ────────────────────────────────────────────────────────────────

func TestGetUserPosts_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	user := &models.User{ID: 1}
	posts := []models.Post{{ID: 1, UserID: 1}, {ID: 2, UserID: 1}}

	userRepo.On("GetUserByID", uint(1)).Return(user, nil)
	postRepo.On("GetPostsByUserID", uint(1), 1, 10).Return(posts, nil)
	postRepo.On("CountPostsByUserID", uint(1)).Return(2, nil)

	resp, err := svc.GetUserPosts(1, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, 2, resp.Total)

	userRepo.AssertExpectations(t)
	postRepo.AssertExpectations(t)
}

func TestGetUserPosts_UserNotFound(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	postRepo := &mocks.MockPostRepository{}
	svc := services.NewPostService(userRepo, postRepo)

	userRepo.On("GetUserByID", uint(99)).Return(nil, fmt.Errorf("not found"))

	resp, err := svc.GetUserPosts(99, 1, 10)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NOT_FOUND")

	userRepo.AssertExpectations(t)
}
