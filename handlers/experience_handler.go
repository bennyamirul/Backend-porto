package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bennyamirul/portfolio-backend/config"
	"github.com/bennyamirul/portfolio-backend/models"
)

// GetExperiences mengambil semua experience dari database dan mengirimkannya dalam response JSON.
// Data diurutkan berdasarkan StartDate terbaru agar timeline portfolio tampil dari pengalaman paling baru.
// Handler menerima context Gin sebagai parameter, lalu menulis response melalui c.JSON tanpa return eksplisit.
func GetExperiences(c *gin.Context) {
	var experiences []models.Experience

	// start_date DESC mengikuti nama kolom GORM untuk field StartDate dan menaruh pengalaman terbaru di atas.
	if err := config.DB.Order("start_date DESC, id ASC").Find(&experiences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data experience"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": experiences})
}
