package routes

import (
	"awesomeProject/handlers"
	"awesomeProject/middleware"
	"awesomeProject/repositories"
	"awesomeProject/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	// Wire up real repository implementations
	userRepo := &repositories.UserRepositoryImpl{}
	tokenRepo := &repositories.RefreshTokenRepositoryImpl{}
	postRepo := &repositories.PostRepositoryImpl{}

	// Wire up services with their dependencies
	authSvc := services.NewAuthService(userRepo, tokenRepo)
	userSvc := services.NewUserService(userRepo)
	postSvc := services.NewPostService(userRepo, postRepo)

	// Wire up handlers with their services
	authH := handlers.NewAuthHandler(authSvc)
	userH := handlers.NewUserHandler(userSvc)
	postH := handlers.NewPostHandler(postSvc)

	r := gin.New()

	// Middlewares globais
	r.Use(middleware.StructuredLogging())
	r.Use(middleware.RateLimiter())

	auth := r.Group("/auth")
	auth.Use(middleware.StrictRateLimiter())
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
		auth.POST("/logout", middleware.RequireAuthWithValidator(authSvc), authH.Logout)
	}

	users := r.Group("/users")
	{
		users.GET("/", userH.GetAllUsers)
		users.GET("/:id/posts", postH.GetUserPosts)
		users.GET("/:id", userH.GetUserByID)

		users.PUT("/:id", middleware.RequireAuthWithValidator(authSvc), userH.UpdateUser)
		users.DELETE("/:id", middleware.RequireAuthWithValidator(authSvc), userH.DeleteUser)
	}

	posts := r.Group("/posts")
	{
		posts.GET("/", postH.GetAllPosts)
		posts.GET("/:id", postH.GetPostByID)

		posts.POST("/", middleware.RequireAuthWithValidator(authSvc), postH.CreatePost)
		posts.PUT("/:id", middleware.RequireAuthWithValidator(authSvc), postH.UpdatePost)
		posts.DELETE("/:id", middleware.RequireAuthWithValidator(authSvc), postH.DeletePost)
	}

	return r
}
