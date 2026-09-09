package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/bennyamirul/portfolio-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase membuka dan menguji koneksi PostgreSQL menggunakan DATABASE_URL.
func ConnectDatabase() error {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		return fmt.Errorf("DATABASE_URL tidak ditemukan")
	}

	dsn, err := withSSLMode(dsn)
	if err != nil {
		return fmt.Errorf("DATABASE_URL tidak valid: %w", err)
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("gagal membuka koneksi PostgreSQL: %w", err)
	}

	DB = database

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan koneksi SQL PostgreSQL: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("gagal melakukan ping PostgreSQL (periksa host, kredensial, jaringan, dan sslmode=require): %w", err)
	}

	// AutoMigrate membuat atau memperbarui tabel berdasarkan struct model tanpa menghapus data yang sudah ada.
	// Setelah aplikasi dijalankan, verifikasi tabel dari terminal PostgreSQL dengan perintah: psql "$DATABASE_URL" -c "\dt"
	if err := DB.AutoMigrate(&models.Project{}, &models.Experience{}, &models.Message{}, &models.Certification{}, &models.Skill{}, &models.Profile{}); err != nil {
		return fmt.Errorf("gagal menjalankan migrasi database: %w", err)
	}

	log.Println("Berhasil terhubung ke PostgreSQL")
	return nil
}

func withSSLMode(dsn string) (string, error) {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return "", fmt.Errorf("scheme %q bukan postgres atau postgresql", parsed.Scheme)
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("host database tidak ditemukan")
	}

	query := parsed.Query()
	if query.Get("sslmode") == "" {
		query.Set("sslmode", "require")
		parsed.RawQuery = query.Encode()
	}

	return parsed.String(), nil
}
