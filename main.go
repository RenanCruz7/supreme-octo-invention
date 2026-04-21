package main

import (
	"fmt"
	"log"

	"awesomeProject/config"
	"awesomeProject/db"
	"awesomeProject/routes"
)

func main() {
	config.Init()

	if err := db.Init(); err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}

	fmt.Println("✓ Banco de dados inicializado com sucesso")

	router := routes.SetupRoutes()

	fmt.Printf("🚀 Servidor Gin iniciado em http://localhost%s\n", config.AppConfig.Port)
	if err := router.Run(config.AppConfig.Port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
