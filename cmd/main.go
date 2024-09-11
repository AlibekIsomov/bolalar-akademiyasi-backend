package main

import (
	"bolalar-akademiyasi/config"
	"bolalar-akademiyasi/database"
	"bolalar-akademiyasi/models"
	"bolalar-akademiyasi/routes"
	_ "github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDatabase(cfg)

	// Auto-migrate the Client model
	database.DB.AutoMigrate(&models.Client{})

	// Setup routes
	router := routes.SetupRouter()

	// Run the server
	router.Run(":8081")
}
