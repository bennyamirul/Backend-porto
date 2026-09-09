package models

import "gorm.io/gorm"

// Message menyimpan pesan dari form contact agar bisa dibaca atau diproses dari backend.
// Embed gorm.Model dipakai supaya GORM otomatis menyediakan ID, CreatedAt, UpdatedAt, dan DeletedAt.
type Message struct {
	gorm.Model

	// Name adalah nama pengirim dan wajib diisi agar pesan memiliki identitas dasar.
	Name string `gorm:"not null"`

	// Email adalah alamat balasan pengirim dan wajib diisi supaya pemilik portfolio bisa menghubungi kembali.
	Email string `gorm:"not null"`

	// Message adalah isi pesan dari form contact, memakai type:text karena panjang pesan pengguna tidak selalu pendek.
	Message string `gorm:"type:text;not null"`
}
