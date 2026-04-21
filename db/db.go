package db

import (
	"awesomeProject/config"
	"awesomeProject/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.AppConfig.DatabasePath), &gorm.Config{})
	if err != nil {
		return err
	}

	return DB.AutoMigrate(&models.User{}, &models.Post{})
}
