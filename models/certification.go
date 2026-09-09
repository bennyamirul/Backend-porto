package models

import "gorm.io/gorm"

// Certification menyimpan data sertifikat yang akan ditampilkan di halaman portfolio.
// Embed gorm.Model dipakai supaya GORM otomatis menyediakan ID, CreatedAt, UpdatedAt, dan DeletedAt.
type Certification struct {
	gorm.Model

	// Title adalah nama sertifikat dan wajib diisi karena menjadi identitas utama certification.
	Title string `gorm:"not null"`

	// Issuer adalah penerbit sertifikat dan wajib diisi agar pengunjung tahu sumber kredensialnya.
	Issuer string `gorm:"not null"`

	// IssueDate memakai string agar format tampilannya fleksibel, misalnya "September 2025".
	IssueDate string

	// CredentialURL menyimpan link ke sertifikat asli dan dibuat opsional karena tidak semua sertifikat punya URL publik.
	CredentialURL string

	// ImageURL menyimpan logo penerbit sertifikat dan dibuat opsional agar data tetap bisa dibuat tanpa gambar.
	ImageURL string
}
