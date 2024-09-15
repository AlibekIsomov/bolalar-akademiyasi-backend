package routes

import (
	"bolalar-akademiyasi/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://127.0.0.1:5501"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	router.GET("/clients", controllers.GetClients)
	router.GET("/clients/", controllers.GetClients)

	router.GET("/clients/:id", controllers.GetClient)
	router.POST("/clients", controllers.CreateClient)
	router.PUT("/clients/:id", controllers.UpdateClient)
	router.DELETE("/clients/:id", controllers.DeleteClient)

	return router
}
