package routes

import (
	"awesomeProject/handlers"
	"awesomeProject/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	users := r.Group("/users")
	{
		users.GET("/", handlers.GetAllUsers)
		users.GET("/:id/posts", handlers.GetUserPosts)
		users.GET("/:id", handlers.GetUserByID)

		users.PUT("/:id", middleware.RequireAuth(), handlers.UpdateUser)
		users.DELETE("/:id", middleware.RequireAuth(), handlers.DeleteUser)
	}

	posts := r.Group("/posts")
	{
		posts.GET("/", handlers.GetAllPosts)
		posts.GET("/:id", handlers.GetPostByID)

		posts.POST("/", middleware.RequireAuth(), handlers.CreatePost)
		posts.PUT("/:id", middleware.RequireAuth(), handlers.UpdatePost)
		posts.DELETE("/:id", middleware.RequireAuth(), handlers.DeletePost)
	}

	return r
}
