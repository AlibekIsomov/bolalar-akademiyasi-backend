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
		AllowOrigins:     []string{"http://185.250.44.239", "http://185.250.44.239/admin", "https://admin.bolalar-akademiyasi.uz", "https://bolalar-akademiyasi.uz", "http://127.0.0.1:5501"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Public routes
	router.POST("/api/auth/login", controllers.Auth)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/api/clients", controllers.GetClients)
		protected.GET("/api/clients/:id", controllers.GetClient)

		protected.POST("/api/clients", controllers.GetClient)

		protected.PUT("/api/clients/:id", controllers.UpdateClient)
		protected.DELETE("/api/clients/:id", controllers.DeleteClient)
	}
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Route not found"})
	})

	return router
}
