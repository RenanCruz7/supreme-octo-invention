package main

import (
	"fmt"
	"log"

	"awesomeProject/db"
	"awesomeProject/routes"
)

func main() {
	// Inicializar banco de dados
	if err := db.Init(); err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}

	fmt.Println("✓ Banco de dados inicializado com sucesso")

	// Configurar rotas
	router := routes.SetupRoutes()

	// Iniciar servidor
	port := ":8080"
	fmt.Printf("🚀 Servidor Gin iniciado em http://localhost%s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
