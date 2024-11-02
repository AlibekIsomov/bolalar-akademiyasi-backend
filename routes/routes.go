package routes

import (
	"bolalar-akademiyasi/controllers"
	"bolalar-akademiyasi/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://185.250.44.239", "http://185.250.44.239/admin", "http://127.0.0.1:5501", "http://127.0.0.1:5502"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Public routes
	router.POST("/auth/login", controllers.Auth)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/clients", controllers.GetClients)
		protected.GET("/clients/:id", controllers.GetClient)

		protected.PUT("/clients/:id", controllers.UpdateClient)
		protected.DELETE("/clients/:id", controllers.DeleteClient)
	}
	router.POST("/clients", controllers.CreateClient)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Route not found"})
	})

	return router
}
