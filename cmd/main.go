package main

import (
	"bolalar-akademiyasi/config"
	"bolalar-akademiyasi/database"
	"bolalar-akademiyasi/models"
	"bolalar-akademiyasi/routes"
	"bolalar-akademiyasi/telegramBot"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Load .env file
	err := godotenv.Load("application-properties.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()

	// Initialize database
	database.InitDatabase(cfg)

	// Auto-migrate the Client model
	database.DB.AutoMigrate(&models.Client{})

	go func() {
		log.Println("Starting Telegram bot...")
		telegramBot.Telegrambot() // Run your telegram bot
	}()

	// Setup routes
	router := routes.SetupRouter()

	// Run the server
	router.Run(":8081")
}
