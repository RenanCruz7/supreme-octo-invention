package repositories

import "awesomeProject/models"

type IUserRepository interface {
	CreateUser(user models.User) (uint, error)
	GetAllUsers(page, limit int) ([]models.User, error)
	CountAllUsers() (int, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user models.User) error
	DeleteUser(id uint) error
}

type IPostRepository interface {
	CreatePost(post models.Post) (uint, error)
	GetAllPosts(page, limit int) ([]models.Post, error)
	CountAllPosts() (int, error)
	GetPostByID(id uint) (*models.Post, error)
	GetPostsByUserID(userID uint, page, limit int) ([]models.Post, error)
	CountPostsByUserID(userID uint) (int, error)
	UpdatePost(post models.Post) error
	DeletePost(id uint) error
}

type IRefreshTokenRepository interface {
	CreateRefreshToken(token models.RefreshToken) error
	GetRefreshToken(token string) (*models.RefreshToken, error)
	RevokeRefreshToken(token string) error
	RevokeAllUserTokens(userID uint) error
}
