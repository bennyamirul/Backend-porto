package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var allowedImageExtensions = map[string]bool{
	".gif":  true,
	".jpeg": true,
	".jpg":  true,
	".png":  true,
	".webp": true,
}

// UploadAdminImage menyimpan satu gambar dari form multipart dan mengembalikan URL publiknya.
func UploadAdminImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file gambar wajib dipilih"})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ukuran gambar maksimal 5 MB"})
		return
	}

	extension := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExtensions[extension] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format gambar harus JPG, PNG, WEBP, atau GIF"})
		return
	}

	// Nama berbasis nanosecond meminimalkan kemungkinan bentrok tanpa dependency tambahan.
	filename := uploadFilename(extension)
	if err := c.SaveUploadedFile(file, filepath.Join("uploads", filename)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan gambar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"imageUrl": "/uploads/" + filename}})
}

func uploadFilename(extension string) string {
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
}