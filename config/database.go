package config

import (
	"log"
	"os"

	"github.com/bennyamirul/portfolio-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase membuka koneksi PostgreSQL menggunakan DATABASE_URL dari environment.
// GORM dipakai di sini karena project sudah memakai GORM, sehingga koneksi dan migrasi schema bisa dikelola dari satu tempat.
func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL tidak ditemukan")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}

	DB = database

	// AutoMigrate membuat atau memperbarui tabel berdasarkan struct model tanpa menghapus data yang sudah ada.
	// Setelah aplikasi dijalankan, verifikasi tabel dari terminal PostgreSQL dengan perintah: psql "$DATABASE_URL" -c "\dt"
	if err := DB.AutoMigrate(&models.Project{}, &models.Experience{}, &models.Message{}, &models.Certification{}, &models.Skill{}, &models.Profile{}); err != nil {
		log.Fatal("Gagal menjalankan migrasi database: ", err)
	}

	log.Println("Berhasil terhubung ke PostgreSQL")
}
