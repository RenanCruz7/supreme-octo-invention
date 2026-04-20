package db

import (
	"awesomeProject/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("./app.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// AutoMigrate cria a tabela automaticamente baseada no modelo
	return DB.AutoMigrate(&models.User{})
}
