package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/bennyamirul/portfolio-backend/config"
	"github.com/bennyamirul/portfolio-backend/models"
)

// GetProjects mengambil semua project dari database dan mengirimkannya dalam response JSON.
// Project diurutkan dengan Featured=true lebih dulu agar frontend bisa langsung menampilkan karya unggulan
// tanpa perlu mengurutkan ulang di browser. Parameter masuk diambil dari context Gin, dan return dikirim
// melalui c.JSON karena handler Gin tidak mengembalikan nilai secara langsung.
func GetProjects(c *gin.Context) {
	var projects []models.Project

	// Featured DESC membuat nilai true muncul sebelum false; id ASC menjaga urutan tetap stabil.
	if err := config.DB.Order("featured DESC, id ASC").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": projects})
}

// GetProjectBySlug mengambil satu project berdasarkan slug dari path URL dan mengirimkannya sebagai JSON.
// Slug dipakai sebagai identifier publik karena lebih ramah URL dibanding ID database. Parameter slug dibaca
// dari c.Param("slug"), dan hasil atau error dikirim lewat c.JSON sesuai pola response API ini.
func GetProjectBySlug(c *gin.Context) {
	var project models.Project
	slug := c.Param("slug")

	// First mengembalikan ErrRecordNotFound saat slug tidak ada, sehingga bisa dibedakan dari error database lain.
	if err := config.DB.Where("slug = ?", slug).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project tidak ditemukan"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})
}
