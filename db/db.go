package db

import (
	"fmt"

	"awesomeProject/config"
	"awesomeProject/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var dialector gorm.Dialector
	var err error

	switch config.AppConfig.DBDriver {
	case "postgres":
		dialector = postgres.Open(buildPostgresConnString())
	case "mysql":
		dialector = mysql.Open(buildMysqlConnString())
	case "sqlite":
		dialector = sqlite.Open(config.AppConfig.DBPath)
	default:
		return fmt.Errorf("driver de banco de dados não suportado: %s", config.AppConfig.DBDriver)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("erro ao conectar banco de dados: %w", err)
	}

	// Auto migrate
	if err := DB.AutoMigrate(&models.User{}, &models.Post{}, &models.RefreshToken{}); err != nil {
		return fmt.Errorf("erro ao executar migrations: %w", err)
	}

	// Configurar connection pool para produção
	if config.AppConfig.DBDriver != "sqlite" {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("erro ao obter banco de dados: %w", err)
		}
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(5)
	}

	return nil
}

func buildPostgresConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
	)
}

func buildMysqlConnString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)
}
