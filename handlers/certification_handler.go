package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bennyamirul/portfolio-backend/config"
	"github.com/bennyamirul/portfolio-backend/models"
)

// GetCertifications mengambil semua certification dari database dan mengirimkannya dalam response JSON.
// IssueDate diurutkan menurun agar certification terbaru tampil lebih dulu, lalu CreatedAt dipakai sebagai
// penentu urutan tambahan jika tanggal tampilannya sama atau kosong. Handler mengirim hasil melalui c.JSON.
func GetCertifications(c *gin.Context) {
	var certifications []models.Certification

	// issue_date adalah nama kolom GORM untuk field IssueDate; created_at DESC menjaga data baru tetap di atas.
	if err := config.DB.Order("issue_date DESC, created_at DESC").Find(&certifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data certification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": certifications})
}
