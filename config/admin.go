package config

import (
	"log"
	"os"
)

// LoadAdminConfig memastikan environment variable admin sudah tersedia sebelum server berjalan.
// ADMIN_PASSWORD adalah password plain text untuk login sederhana di development, sedangkan
// ADMIN_JWT_SECRET adalah kunci rahasia untuk menandatangani JWT agar token tidak bisa dipalsukan.
// Fungsi ini tidak menerima parameter dan tidak mengembalikan nilai; kalau konfigurasi wajib kosong,
// aplikasi dihentikan lebih awal dengan log.Fatal supaya error konfigurasi terlihat jelas.
func LoadAdminConfig() {
	if os.Getenv("ADMIN_PASSWORD") == "" {
		log.Fatal("ADMIN_PASSWORD tidak ditemukan")
	}

	if os.Getenv("ADMIN_JWT_SECRET") == "" {
		log.Fatal("ADMIN_JWT_SECRET tidak ditemukan")
	}
}
