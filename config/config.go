package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabasePath string
	Port         string
	JWTSecret    string
}

var AppConfig *Config

func Init() {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	AppConfig = &Config{
		DatabasePath: getEnv("DATABASE_PATH", "./app.db"),
		Port:         port,
		JWTSecret:    getEnv("JWT_SECRET", "default_secret_key"),
	}

	if AppConfig.JWTSecret == "default_secret_key" {
		log.Println("⚠️  Aviso: Usando JWT_SECRET padrão. Configure JWT_SECRET no arquivo .env para produção!")
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
