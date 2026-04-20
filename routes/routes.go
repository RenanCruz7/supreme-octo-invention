package routes

import (
	"awesomeProject/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	users := r.Group("/users")
	{
		users.POST("/", handlers.CreateUser)
		users.GET("/", handlers.GetAllUsers)
		users.GET("/:id", handlers.GetUserByID)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
		users.GET("/:userId/posts", handlers.GetUserPosts) // Posts de um usuário
	}

	posts := r.Group("/posts")
	{
		posts.POST("/", handlers.CreatePost)
		posts.GET("/", handlers.GetAllPosts)
		posts.GET("/:id", handlers.GetPostByID)
		posts.PUT("/:id", handlers.UpdatePost)
		posts.DELETE("/:id", handlers.DeletePost)
	}

	return r
}
