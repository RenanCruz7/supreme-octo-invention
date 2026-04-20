package repositories

import (
	"awesomeProject/db"
	"awesomeProject/models"
)

// CreateUser insere um novo usuário no banco de dados
func CreateUser(user models.User) (uint, error) {
	result := db.DB.Create(&user)
	return user.ID, result.Error
}

// GetAllUsers retorna todos os usuários
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := db.DB.Find(&users)
	return users, result.Error
}

// GetUserByID retorna um usuário pelo ID
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := db.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser atualiza um usuário existente
func UpdateUser(user models.User) error {
	return db.DB.Save(&user).Error
}

// DeleteUser remove um usuário
func DeleteUser(id uint) error {
	return db.DB.Delete(&models.User{}, id).Error
}
