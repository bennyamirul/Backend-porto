package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/bennyamirul/portfolio-backend/config"
	"github.com/bennyamirul/portfolio-backend/models"
)

type skillItemRequest struct {
    Name     string `json:"name"`
    Category string `json:"category"`
}

type skillRequest struct {
    Name     string            `json:"name"`
    Category string            `json:"category"`
    Skills   []skillItemRequest `json:"skills"`
}

// GetSkills mengambil semua skill dari database lalu mengirimkannya sebagai array flat.
// Frontend bisa mengelompokkan skill menjadi marquee dan grid kategori sesuai kebutuhan UI.
func GetSkills(c *gin.Context) {
    var skills []models.Skill

    if err := config.DB.Order("category ASC, id ASC").Find(&skills).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": skills})
}

// GetAdminSkills mengambil semua skill untuk dashboard admin dalam urutan yang mudah diatur.
func GetAdminSkills(c *gin.Context) {
    GetSkills(c)
}

// CreateAdminSkill membuat satu atau beberapa skill baru dari payload admin.
// Format lama yang masih didukung: {"name":"React","category":"Frontend"}
// Format baru: {"skills":[{"name":"React","category":"Frontend"},{"name":"Go","category":"Backend"}]}
func CreateAdminSkill(c *gin.Context) {
    var request skillRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
        return
    }

    if len(request.Skills) > 0 {
        createdSkills := make([]models.Skill, 0, len(request.Skills))

        for _, item := range request.Skills {
            name := strings.TrimSpace(item.Name)
            category := strings.TrimSpace(item.Category)
            if name == "" || category == "" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "setiap skill harus punya name dan category"})
                return
            }

            createdSkills = append(createdSkills, models.Skill{
                Name:     name,
                Category: category,
            })
        }

        if err := config.DB.Create(&createdSkills).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat skill"})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"data": createdSkills})
        return
    }

    name := strings.TrimSpace(request.Name)
    category := strings.TrimSpace(request.Category)
    if name == "" || category == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "name dan category wajib diisi"})
        return
    }

    skill := models.Skill{
        Name:     name,
        Category: category,
    }

    if err := config.DB.Create(&skill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat skill"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"data": skill})
}

// UpdateAdminSkill memperbarui skill yang sudah ada berdasarkan ID.
func UpdateAdminSkill(c *gin.Context) {
    id, ok := parseIDParam(c)
    if !ok {
        return
    }

    var request skillRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
        return
    }

    name := strings.TrimSpace(request.Name)
    category := strings.TrimSpace(request.Category)
    if name == "" || category == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "name dan category wajib diisi"})
        return
    }

    var skill models.Skill
    if err := config.DB.First(&skill, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "skill tidak ditemukan"})
            return
        }

        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil skill"})
        return
    }

    skill.Name = name
    skill.Category = category

    if err := config.DB.Save(&skill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memperbarui skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": skill})
}

// DeleteAdminSkill menghapus skill berdasarkan ID dengan soft delete GORM.
func DeleteAdminSkill(c *gin.Context) {
    id, ok := parseIDParam(c)
    if !ok {
        return
    }

    if err := config.DB.Delete(&models.Skill{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Skill berhasil dihapus"})
}
