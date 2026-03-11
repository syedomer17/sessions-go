package routes

import (
	"session-demo/internal/controllers"
	"session-demo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
	
		auth.POST("/login", controllers.Login)
	}

	protected := r.Group("/api")

	protected.Use(middleware.SessionMiddleware())
	{
		protected.GET("/profile", controllers.Profile)
	}

	return r
}
