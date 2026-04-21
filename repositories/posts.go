package repositories

import (
	"awesomeProject/db"
	"awesomeProject/models"
)

func CreatePost(post models.Post) (uint, error) {
	result := db.DB.Create(&post)
	return post.ID, result.Error
}

func GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	result := db.DB.Preload("User").Find(&posts)
	return posts, result.Error
}

func GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	result := db.DB.Preload("User").First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

func GetPostsByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	result := db.DB.Preload("User").Where("user_id = ?", userID).Find(&posts)
	return posts, result.Error
}

func UpdatePost(post models.Post) error {
	return db.DB.Save(&post).Error
}

func DeletePost(id uint) error {
	return db.DB.Delete(&models.Post{}, id).Error
}
