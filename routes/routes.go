package routes

import (
	"awesomeProject/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Rotas de usuários
	users := r.Group("/users")
	{
		users.POST("/", handlers.CreateUser)
		users.GET("/", handlers.GetAllUsers)
		users.GET("/:id", handlers.GetUserByID)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}

	return r
}
