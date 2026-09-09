package models

import "gorm.io/gorm"

// Project menyimpan data karya/portfolio yang akan ditampilkan di website.
// Embed gorm.Model dipakai supaya GORM otomatis menyediakan ID, CreatedAt, UpdatedAt, dan DeletedAt.
type Project struct {
	gorm.Model

	// Title adalah nama project yang tampil di halaman portfolio, jadi wajib diisi dengan tag not null.
	Title string `gorm:"not null"`

	// Slug adalah versi URL-friendly dari title dan harus unik agar route detail project tidak ambigu.
	Slug string `gorm:"uniqueIndex;not null"`

	// Description adalah ringkasan pendek project, disimpan sebagai text agar tidak cepat mentok batas panjang.
	Description string `gorm:"type:text;not null"`

	// Content adalah isi/detail panjang project, sehingga memakai type:text agar PostgreSQL menyimpannya sebagai teks panjang.
	Content string `gorm:"type:text;not null"`

	// TechStack menyimpan daftar teknologi dalam bentuk teks agar sederhana untuk tahap awal backend.
	TechStack string `gorm:"type:text;not null"`

	// ImageURL menyimpan alamat gambar project dan dibuat opsional karena tidak semua project harus punya gambar.
	ImageURL string

	// RepoURL menyimpan alamat repository source code dan dibuat opsional karena project tertentu bisa private.
	RepoURL string

	// DemoURL menyimpan alamat demo live dan dibuat opsional karena tidak semua project punya deployment publik.
	DemoURL string

	// Featured menandai project unggulan, dengan default false supaya project baru tidak otomatis muncul sebagai highlight.
	Featured bool `gorm:"default:false;not null"`
}
