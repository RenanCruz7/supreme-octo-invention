package main

import (
	"fmt"
	"log"
	"time"

	"awesomeProject/config"
	"awesomeProject/db"
	"awesomeProject/repositories"
	"awesomeProject/routes"
)

func main() {
	config.Init()

	if err := db.Init(); err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}

	fmt.Println("✓ Banco de dados inicializado com sucesso")

	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := repositories.DeleteExpiredTokens(); err != nil {
				log.Printf("⚠️  Erro ao limpar tokens expirados: %v", err)
			}
		}
	}()

	router := routes.SetupRoutes()

	fmt.Printf("🚀 Servidor Gin iniciado em http://localhost%s\n", config.AppConfig.Port)
	if err := router.Run(config.AppConfig.Port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
