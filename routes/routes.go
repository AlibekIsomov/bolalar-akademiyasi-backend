package routes

import (
	"bolalar-akademiyasi/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Apply CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://127.0.0.1:5500"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	// Define routes
	clientRoutes := router.Group("/clients")
	{
		clientRoutes.GET("/", controllers.GetClients)
		clientRoutes.GET("/:id", controllers.GetClient)
		clientRoutes.POST("/", controllers.CreateClient)
		clientRoutes.PUT("/:id", controllers.UpdateClient)
		clientRoutes.DELETE("/:id", controllers.DeleteClient)
	}

	return router
}
