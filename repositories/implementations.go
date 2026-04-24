package repositories

import "awesomeProject/models"

// UserRepositoryImpl adapts package-level functions to IUserRepository.
type UserRepositoryImpl struct{}

func (r *UserRepositoryImpl) CreateUser(user models.User) (uint, error) {
	return CreateUser(user)
}
func (r *UserRepositoryImpl) GetAllUsers(page, limit int) ([]models.User, error) {
	return GetAllUsers(page, limit)
}
func (r *UserRepositoryImpl) CountAllUsers() (int, error) { return CountAllUsers() }
func (r *UserRepositoryImpl) GetUserByID(id uint) (*models.User, error) {
	return GetUserByID(id)
}
func (r *UserRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	return GetUserByEmail(email)
}
func (r *UserRepositoryImpl) UpdateUser(user models.User) error { return UpdateUser(user) }
func (r *UserRepositoryImpl) DeleteUser(id uint) error          { return DeleteUser(id) }

// PostRepositoryImpl adapts package-level functions to IPostRepository.
type PostRepositoryImpl struct{}

func (r *PostRepositoryImpl) CreatePost(post models.Post) (uint, error) {
	return CreatePost(post)
}
func (r *PostRepositoryImpl) GetAllPosts(page, limit int) ([]models.Post, error) {
	return GetAllPosts(page, limit)
}
func (r *PostRepositoryImpl) CountAllPosts() (int, error) { return CountAllPosts() }
func (r *PostRepositoryImpl) GetPostByID(id uint) (*models.Post, error) {
	return GetPostByID(id)
}
func (r *PostRepositoryImpl) GetPostsByUserID(userID uint, page, limit int) ([]models.Post, error) {
	return GetPostsByUserID(userID, page, limit)
}
func (r *PostRepositoryImpl) CountPostsByUserID(userID uint) (int, error) {
	return CountPostsByUserID(userID)
}
func (r *PostRepositoryImpl) UpdatePost(post models.Post) error { return UpdatePost(post) }
func (r *PostRepositoryImpl) DeletePost(id uint) error          { return DeletePost(id) }

// RefreshTokenRepositoryImpl adapts package-level functions to IRefreshTokenRepository.
type RefreshTokenRepositoryImpl struct{}

func (r *RefreshTokenRepositoryImpl) CreateRefreshToken(token models.RefreshToken) error {
	return CreateRefreshToken(token)
}
func (r *RefreshTokenRepositoryImpl) GetRefreshToken(token string) (*models.RefreshToken, error) {
	return GetRefreshToken(token)
}
func (r *RefreshTokenRepositoryImpl) RevokeRefreshToken(token string) error {
	return RevokeRefreshToken(token)
}
func (r *RefreshTokenRepositoryImpl) RevokeAllUserTokens(userID uint) error {
	return RevokeAllUserTokens(userID)
}
