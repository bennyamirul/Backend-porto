package models

import "gorm.io/gorm"

// Profile menyimpan data profil utama yang ditampilkan di sidebar portofolio.
// Aplikasi memakai satu row profil agar data nama, email, dan foto bisa diubah dari admin panel tanpa
// menambah tabel tambahan atau model terpisah.
type Profile struct {
	gorm.Model

	Name      string `gorm:"not null"`
	Email     string `gorm:"not null;index"`
	AvatarURL string
}
