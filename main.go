package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/bennyamirul/portfolio-backend/config"
	"github.com/bennyamirul/portfolio-backend/router"
)

func main() {
	// Load .env hanya kalau ada — jangan fatal kalau tidak ketemu (misal di production/Railway)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config.LoadAdminConfig()
	config.ConnectDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback untuk local dev
	}

	router.SetupRouter().Run(":" + port)
}