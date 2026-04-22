package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPath     string // Apenas para SQLite

	// Server
	Port      string
	JWTSecret string
}

var AppConfig *Config

func Init() {
	_ = godotenv.Load()

	dbDriver := getEnv("DB_DRIVER", "sqlite")
	if !isValidDriver(dbDriver) {
		log.Fatalf("❌ DB_DRIVER inválido: %s. Use: sqlite, postgres ou mysql", dbDriver)
	}

	port := getEnv("PORT", "8080")
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	AppConfig = &Config{
		// Database
		DBDriver:   dbDriver,
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", getDefaultPort(dbDriver)),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "awesome_db"),
		DBPath:     getEnv("DB_PATH", "./app.db"),

		// Server
		Port:      port,
		JWTSecret: getEnv("JWT_SECRET", "default_secret_key"),
	}

	if AppConfig.JWTSecret == "default_secret_key" {
		log.Println("⚠️  Aviso: Usando JWT_SECRET padrão. Configure JWT_SECRET no arquivo .env para produção!")
	}

	log.Printf("📦 Banco de dados configurado: %s\n", AppConfig.DBDriver)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func isValidDriver(driver string) bool {
	validDrivers := map[string]bool{
		"sqlite":   true,
		"postgres": true,
		"mysql":    true,
	}
	return validDrivers[driver]
}

func getDefaultPort(driver string) string {
	ports := map[string]string{
		"postgres": "5432",
		"mysql":    "3306",
		"sqlite":   "",
	}
	return ports[driver]
}
