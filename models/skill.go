package models

import "gorm.io/gorm"

// Skill menyimpan teknologi atau kemampuan yang ditampilkan di homepage dan di admin.
// Category digunakan untuk mengelompokkan skill di UI agar tidak perlu field urutan tambahan.
type Skill struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Category string `gorm:"not null;index"`
}
