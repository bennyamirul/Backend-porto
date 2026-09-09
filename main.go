package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/bennyamirul/portfolio-backend/config"
	"github.com/bennyamirul/portfolio-backend/router"
)

// main memuat konfigurasi environment, menghubungkan database, lalu menjalankan HTTP server Gin.
// Alur ini dipakai agar aplikasi berhenti lebih awal kalau file .env atau database belum siap.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.LoadAdminConfig()
	config.ConnectDatabase()

	router.SetupRouter().Run(":8080")
}
