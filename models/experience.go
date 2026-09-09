package models

import (
	"time"

	"gorm.io/gorm"
)

// Experience menyimpan riwayat pekerjaan atau pengalaman profesional untuk halaman portfolio.
// Embed gorm.Model dipakai supaya GORM otomatis menyediakan ID, CreatedAt, UpdatedAt, dan DeletedAt.
type Experience struct {
	gorm.Model

	// Company adalah nama perusahaan/organisasi dan wajib diisi karena menjadi identitas utama pengalaman.
	Company string `gorm:"not null"`

	ImageURL string

	Location string

	// Role adalah posisi atau tanggung jawab utama, sehingga wajib diisi untuk konteks pengalaman.
	Role string `gorm:"not null"`

	// StartDate memakai time.Time karena PostgreSQL dan GORM bisa memetakannya ke tipe tanggal/waktu dengan aman.
	StartDate time.Time `gorm:"not null"`

	// EndDate memakai pointer agar nil bisa mewakili pengalaman yang masih aktif atau belum memiliki tanggal selesai.
	EndDate *time.Time

	// Description menjelaskan detail pekerjaan, memakai type:text karena isi pengalaman bisa lebih panjang dari satu kalimat.
	Description string `gorm:"type:text;not null"`
}
