package repositories

import (
	"awesomeProject/db"
	"awesomeProject/models"
)

// CreatePost insere um novo post no banco de dados
func CreatePost(post models.Post) (uint, error) {
	result := db.DB.Create(&post)
	return post.ID, result.Error
}

// GetAllPosts retorna todos os posts
func GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	result := db.DB.Preload("User").Find(&posts) // Carrega o usuário relacionado
	return posts, result.Error
}

// GetPostByID retorna um post pelo ID
func GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	result := db.DB.Preload("User").First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

// GetPostsByUserID retorna todos os posts de um usuário
func GetPostsByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	result := db.DB.Preload("User").Where("user_id = ?", userID).Find(&posts)
	return posts, result.Error
}

// UpdatePost atualiza um post existente
func UpdatePost(post models.Post) error {
	return db.DB.Save(&post).Error
}

// DeletePost remove um post
func DeletePost(id uint) error {
	return db.DB.Delete(&models.Post{}, id).Error
}
