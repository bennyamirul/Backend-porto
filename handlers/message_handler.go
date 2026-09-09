package handlers

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/bennyamirul/portfolio-backend/config"
	"github.com/bennyamirul/portfolio-backend/models"
)

type createMessageRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// CreateMessage menerima payload JSON contact form, memvalidasinya, lalu menyimpan pesan ke database.
// Validasi dilakukan di handler agar response error bisa konsisten dan mudah dipahami oleh frontend.
// Body JSON dibaca dari context Gin, sedangkan hasil sukses atau gagal dikirim lewat c.JSON.
func CreateMessage(c *gin.Context) {
	var request createMessageRequest

	// ShouldBindJSON memastikan body request benar-benar JSON yang bisa dipetakan ke struct request.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	name := strings.TrimSpace(request.Name)
	email := strings.TrimSpace(request.Email)
	messageText := strings.TrimSpace(request.Message)

	// Field kosong ditolak supaya data yang tersimpan selalu cukup lengkap untuk ditindaklanjuti.
	if name == "" || email == "" || messageText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email, dan message wajib diisi"})
		return
	}

	// net/mail dipakai agar validasi email mengikuti parser standar Go, bukan pengecekan string manual.
	if _, err := mail.ParseAddress(email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format email tidak valid"})
		return
	}

	message := models.Message{
		Name:    name,
		Email:   email,
		Message: messageText,
	}

	if err := config.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan pesan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": message})
}
