package services

import "awesomeProject/models"

// IAuthService defines the auth business logic contract used by handlers.
type IAuthService interface {
	Register(req models.RegisterRequest) (*models.AuthResponse, error)
	Login(req models.LoginRequest) (*models.AuthResponse, error)
	RefreshTokens(refreshToken string) (*models.AuthResponse, error)
	Logout(userID uint) error
}

// IUserService defines the user business logic contract used by handlers.
type IUserService interface {
	GetAllUsers(page, limit int) (*models.PaginatedResponse[models.User], error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id uint) error
}

// IPostService defines the post business logic contract used by handlers.
type IPostService interface {
	CreatePost(post models.Post, userID uint) (*models.Post, error)
	GetAllPosts(page, limit int) (*models.PaginatedResponse[models.Post], error)
	GetPostByID(id uint) (*models.Post, error)
	GetUserPosts(userID uint, page, limit int) (*models.PaginatedResponse[models.Post], error)
	UpdatePost(id uint, post models.Post, userID uint) (*models.Post, error)
	DeletePost(id uint, userID uint) error
}
