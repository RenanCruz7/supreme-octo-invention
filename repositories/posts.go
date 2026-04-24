package repositories

import (
	"awesomeProject/db"
	"awesomeProject/models"
)

func CreatePost(post models.Post) (uint, error) {
	result := db.DB.Create(&post)
	return post.ID, result.Error
}

func CountAllPosts() (int, error) {
	var count int64
	result := db.DB.Model(&models.Post{}).Count(&count)
	return int(count), result.Error
}

func GetAllPosts(page, limit int) ([]models.Post, error) {
	var posts []models.Post
	offset := (page - 1) * limit
	result := db.DB.Preload("User").Offset(offset).Limit(limit).Find(&posts)
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

func CountPostsByUserID(userID uint) (int, error) {
	var count int64
	result := db.DB.Model(&models.Post{}).Where("user_id = ?", userID).Count(&count)
	return int(count), result.Error
}

func GetPostsByUserID(userID uint, page, limit int) ([]models.Post, error) {
	var posts []models.Post
	offset := (page - 1) * limit
	result := db.DB.Preload("User").Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&posts)
	return posts, result.Error
}

func UpdatePost(post models.Post) error {
	return db.DB.Save(&post).Error
}

func DeletePost(id uint) error {
	return db.DB.Delete(&models.Post{}, id).Error
}
