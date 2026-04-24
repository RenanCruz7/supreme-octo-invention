package repositories

import (
	"awesomeProject/db"
	"awesomeProject/models"
)

func CreateUser(user models.User) (uint, error) {
	result := db.DB.Create(&user)
	return user.ID, result.Error
}

func CountAllUsers() (int, error) {
	var count int64
	result := db.DB.Model(&models.User{}).Count(&count)
	return int(count), result.Error
}

func GetAllUsers(page, limit int) ([]models.User, error) {
	var users []models.User
	offset := (page - 1) * limit
	result := db.DB.Offset(offset).Limit(limit).Find(&users)
	return users, result.Error
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := db.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UpdateUser(user models.User) error {
	return db.DB.Save(&user).Error
}

func DeleteUser(id uint) error {
	return db.DB.Delete(&models.User{}, id).Error
}
