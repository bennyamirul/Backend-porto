package handlers

import (
	"errors"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/bennyamirul/portfolio-backend/config"
	"github.com/bennyamirul/portfolio-backend/models"
)

type adminLoginRequest struct {
	Password string `json:"password"`
}

type projectRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Content     string `json:"content"`
	TechStack   string `json:"techStack"`
	ImageURL    string `json:"imageUrl"`
	RepoURL     string `json:"repoUrl"`
	DemoURL     string `json:"demoUrl"`
	Featured    bool   `json:"featured"`
}

type experienceRequest struct {
	Company     string     `json:"company"`
	ImageURL    string     `json:"imageUrl"`
	Location    string     `json:"location"`
	Role        string     `json:"role"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Description string     `json:"description"`
}

type certificationRequest struct {
	Title         string `json:"title"`
	Issuer        string `json:"issuer"`
	IssueDate     string `json:"issueDate"`
	CredentialURL string `json:"credentialUrl"`
	ImageURL      string `json:"imageUrl"`
}

type profileRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl"`
}

// AdminLogin memeriksa password admin lalu membuat JWT 24 jam jika password benar.
// Password dibandingkan dengan ADMIN_PASSWORD dari environment karena kebutuhan project masih admin sederhana.
// JWT disimpan sebagai cookie httpOnly bernama admin_token supaya browser mengirimnya otomatis, tetapi JavaScript
// frontend tidak bisa membaca isi token. Parameter body adalah JSON {"password":"..."}, dan response suksesnya
// adalah {"message":"Login berhasil"}; kalau password salah, handler mengembalikan 401.
func AdminLogin(c *gin.Context) {
	var request adminLoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	if request.Password != os.Getenv("ADMIN_PASSWORD") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	claims := jwt.RegisteredClaims{
		Subject:   "admin",
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("ADMIN_JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat token"})
		return
	}

	// Secure=false dipakai untuk localhost HTTP. Di production HTTPS, nilai ini sebaiknya true.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("admin_token", tokenString, int(24*time.Hour/time.Second), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login berhasil"})
}

// AdminLogout menghapus cookie admin_token dengan mengirim cookie kosong yang langsung kedaluwarsa.
// Browser akan mengganti cookie lama berdasarkan nama/path/domain yang sama, sehingga request admin berikutnya
// tidak lagi membawa token. Handler ini tidak membutuhkan body dan mengembalikan pesan sukses sederhana.
func AdminLogout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("admin_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

// CreateAdminProject membuat data project baru dari JSON admin.
// Title dan Slug wajib diisi karena keduanya menjadi identitas utama project, dan Slug dicek unik
// sebelum insert agar URL detail publik tidak bentrok. Parameter body berisi field project, dan return-nya
// adalah data project yang baru dibuat atau error validasi/database.
func CreateAdminProject(c *gin.Context) {
	var request projectRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	title := strings.TrimSpace(request.Title)
	slug := strings.TrimSpace(request.Slug)
	if title == "" || slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title dan slug wajib diisi"})
		return
	}

	if slugExists(slug, 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug sudah dipakai project lain"})
		return
	}

	project := models.Project{
		Title:       title,
		Slug:        slug,
		Description: strings.TrimSpace(request.Description),
		Content:     strings.TrimSpace(request.Content),
		TechStack:   strings.TrimSpace(request.TechStack),
		ImageURL:    strings.TrimSpace(request.ImageURL),
		RepoURL:     strings.TrimSpace(request.RepoURL),
		DemoURL:     strings.TrimSpace(request.DemoURL),
		Featured:    request.Featured,
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": project})
}

// UpdateAdminProject memperbarui project berdasarkan ID path.
// Handler mencari data lama lebih dulu supaya update pada ID yang tidak ada bisa mengembalikan 404,
// lalu mengecek ulang slug agar tidak dipakai oleh project lain. Parameter :id berasal dari URL, body
// berisi field project versi baru, dan return-nya adalah project yang sudah diperbarui.
func UpdateAdminProject(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var request projectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	title := strings.TrimSpace(request.Title)
	slug := strings.TrimSpace(request.Slug)
	if title == "" || slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title dan slug wajib diisi"})
		return
	}

	var project models.Project
	if err := config.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project tidak ditemukan"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil project"})
		return
	}

	if slugExists(slug, id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug sudah dipakai project lain"})
		return
	}

	project.Title = title
	project.Slug = slug
	project.Description = strings.TrimSpace(request.Description)
	project.Content = strings.TrimSpace(request.Content)
	project.TechStack = strings.TrimSpace(request.TechStack)
	project.ImageURL = strings.TrimSpace(request.ImageURL)
	project.RepoURL = strings.TrimSpace(request.RepoURL)
	project.DemoURL = strings.TrimSpace(request.DemoURL)
	project.Featured = request.Featured

	if err := config.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memperbarui project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})
}

// DeleteAdminProject menghapus project berdasarkan ID path dengan soft delete bawaan gorm.Model.
// Soft delete berarti baris database tidak benar-benar hilang; GORM mengisi DeletedAt sehingga query normal
// tidak lagi menampilkan data tersebut. Parameter :id berasal dari URL, dan return suksesnya berupa message.
func DeleteAdminProject(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	if err := config.DB.Delete(&models.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project berhasil dihapus"})
}

// CreateAdminExperience membuat data experience baru dari JSON admin.
// Body dipetakan ke model Experience agar format tanggal tetap ditangani parser JSON Go untuk time.Time.
// Parameter body berisi company, role, startDate, endDate opsional, dan description; return-nya adalah data baru.
func CreateAdminExperience(c *gin.Context) {
	var request experienceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	experience := models.Experience{
		Company:     strings.TrimSpace(request.Company),
		ImageURL:    strings.TrimSpace(request.ImageURL),
		Location:    strings.TrimSpace(request.Location),
		Role:        strings.TrimSpace(request.Role),
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		Description: strings.TrimSpace(request.Description),
	}

	if err := config.DB.Create(&experience).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat experience"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": experience})
}

// UpdateAdminExperience memperbarui experience berdasarkan ID path.
// Data lama dicari lebih dulu supaya endpoint memberi 404 untuk ID yang tidak ada, lalu field dari body
// menggantikan nilai lama. Parameter :id berasal dari URL, body berisi field experience, dan return-nya data terbaru.
func UpdateAdminExperience(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var request experienceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	var experience models.Experience
	if err := config.DB.First(&experience, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "experience tidak ditemukan"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil experience"})
		return
	}

	experience.Company = strings.TrimSpace(request.Company)
	experience.ImageURL = strings.TrimSpace(request.ImageURL)
	experience.Location = strings.TrimSpace(request.Location)
	experience.Role = strings.TrimSpace(request.Role)
	experience.StartDate = request.StartDate
	experience.EndDate = request.EndDate
	experience.Description = strings.TrimSpace(request.Description)

	if err := config.DB.Save(&experience).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memperbarui experience"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": experience})
}

// DeleteAdminExperience menghapus experience berdasarkan ID path dengan soft delete bawaan gorm.Model.
// Karena model embed gorm.Model, GORM mengisi DeletedAt dan menyembunyikan data dari query normal tanpa
// menghapus baris fisiknya. Parameter :id berasal dari URL, dan return suksesnya berupa message.
func DeleteAdminExperience(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	if err := config.DB.Delete(&models.Experience{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus experience"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Experience berhasil dihapus"})
}

// CreateAdminCertification membuat data certification baru dari JSON admin.
// Title dan Issuer wajib diisi karena keduanya menjelaskan sertifikat dan penerbitnya. IssueDate tetap string
// agar admin bebas memakai format tampilan seperti "September 2025"; return-nya adalah data baru atau error.
func CreateAdminCertification(c *gin.Context) {
	var request certificationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	title := strings.TrimSpace(request.Title)
	issuer := strings.TrimSpace(request.Issuer)
	if title == "" || issuer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title dan issuer wajib diisi"})
		return
	}

	certification := models.Certification{
		Title:         title,
		Issuer:        issuer,
		IssueDate:     strings.TrimSpace(request.IssueDate),
		CredentialURL: strings.TrimSpace(request.CredentialURL),
		ImageURL:      strings.TrimSpace(request.ImageURL),
	}

	if err := config.DB.Create(&certification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat certification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": certification})
}

// UpdateAdminCertification memperbarui certification berdasarkan ID path.
// Data lama dicari lebih dulu agar ID yang tidak ada mendapat response 404. Field body kemudian menimpa data
// lama setelah trim spasi supaya data yang tersimpan lebih rapi dan konsisten dengan handler admin lain.
func UpdateAdminCertification(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var request certificationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	title := strings.TrimSpace(request.Title)
	issuer := strings.TrimSpace(request.Issuer)
	if title == "" || issuer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title dan issuer wajib diisi"})
		return
	}

	var certification models.Certification
	if err := config.DB.First(&certification, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "certification tidak ditemukan"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil certification"})
		return
	}

	certification.Title = title
	certification.Issuer = issuer
	certification.IssueDate = strings.TrimSpace(request.IssueDate)
	certification.CredentialURL = strings.TrimSpace(request.CredentialURL)
	certification.ImageURL = strings.TrimSpace(request.ImageURL)

	if err := config.DB.Save(&certification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memperbarui certification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": certification})
}

// DeleteAdminCertification menghapus certification berdasarkan ID path dengan soft delete bawaan gorm.Model.
// Soft delete mengisi DeletedAt, bukan menghapus baris fisik, sehingga data masih bisa diaudit atau dipulihkan
// melalui query Unscoped jika suatu saat diperlukan. Return suksesnya berupa message sederhana.
func DeleteAdminCertification(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	if err := config.DB.Delete(&models.Certification{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus certification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certification berhasil dihapus"})
}

// GetAdminMessages mengambil semua pesan contact untuk dashboard admin.
// Pesan diurutkan dari CreatedAt terbaru agar admin melihat pesan terbaru lebih dulu. Handler tidak menerima
// body; return-nya adalah array Message dalam key data atau error database.
func GetAdminMessages(c *gin.Context) {
	var messages []models.Message

	if err := config.DB.Order("created_at DESC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil pesan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": messages})
}

func GetProfile(c *gin.Context) {
	var profile models.Profile

	if err := config.DB.Order("id ASC").First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			profile = models.Profile{
				Name:  "Benny Amirul",
				Email: "bennyamirul@gmail.com",
			}
			if err := config.DB.Create(&profile).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat profil default"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": profile})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile})
}

func GetAdminProfile(c *gin.Context) {
	GetProfile(c)
}

func UpdateAdminProfile(c *gin.Context) {
	var request profileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format JSON tidak valid"})
		return
	}

	name := strings.TrimSpace(request.Name)
	email := strings.TrimSpace(request.Email)
	if name == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name dan email wajib diisi"})
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format email tidak valid"})
		return
	}

	var profile models.Profile
	if err := config.DB.Order("id ASC").First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			profile = models.Profile{
				Name:      name,
				Email:     email,
				AvatarURL: strings.TrimSpace(request.AvatarURL),
			}
			if err := config.DB.Create(&profile).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat profil"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"data": profile})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil profil"})
		return
	}

	profile.Name = name
	profile.Email = email
	profile.AvatarURL = strings.TrimSpace(request.AvatarURL)

	if err := config.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memperbarui profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile})
}

// parseIDParam membaca parameter :id dari URL dan mengubahnya menjadi uint untuk query GORM.
// Gin menyimpan parameter path sebagai string, sedangkan primary key gorm.Model bertipe uint, jadi helper ini
// menjaga validasi ID tetap konsisten. Parameter masuknya adalah context Gin; return-nya ID dan boolean sukses.
func parseIDParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return 0, false
	}

	return uint(id), true
}

// slugExists mengecek apakah slug sudah dipakai project lain.
// exceptID dipakai saat update agar project yang sedang diedit boleh mempertahankan slug miliknya sendiri.
// Parameter slug adalah slug yang dicek, exceptID adalah ID yang dikecualikan, dan return true berarti bentrok.
func slugExists(slug string, exceptID uint) bool {
	var count int64
	query := config.DB.Unscoped().Model(&models.Project{}).Where("slug = ?", slug)
	if exceptID != 0 {
		query = query.Where("id <> ?", exceptID)
	}

	return query.Count(&count).Error == nil && count > 0
}
