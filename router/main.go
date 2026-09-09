package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bennyamirul/portfolio-backend/handlers"
	"github.com/bennyamirul/portfolio-backend/middleware"
)

// SetupRouter membuat router Gin, memasang middleware CORS, dan mendaftarkan semua endpoint API.
// Fungsi ini dipisah dari main agar konfigurasi route terkumpul di satu tempat dan mudah dites/diubah.
// Return-nya adalah *gin.Engine yang kemudian dijalankan oleh main.go sebagai HTTP server.
func SetupRouter() *gin.Engine {
	router := gin.Default()
	_ = os.MkdirAll("uploads", 0755)
	router.Static("/uploads", "./uploads")

	router.Use(corsMiddleware())

	// Route root ini dipakai sebagai health check sederhana untuk memastikan server Gin sudah hidup.
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Portfolio API is running!",
		})
	})

	api := router.Group("/api")
	{
		api.GET("/projects", handlers.GetProjects)
		api.GET("/projects/:slug", handlers.GetProjectBySlug)
		api.GET("/experience", handlers.GetExperiences)
		api.GET("/certifications", handlers.GetCertifications)
		api.GET("/skills", handlers.GetSkills)
		api.GET("/profile", handlers.GetProfile)
		api.POST("/messages", handlers.CreateMessage)

		api.POST("/admin/login", handlers.AdminLogin)
		api.POST("/admin/logout", handlers.AdminLogout)

		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware())
		{
			admin.POST("/upload", handlers.UploadAdminImage)
			admin.POST("/projects", handlers.CreateAdminProject)
			admin.PUT("/projects/:id", handlers.UpdateAdminProject)
			admin.DELETE("/projects/:id", handlers.DeleteAdminProject)
			admin.POST("/experience", handlers.CreateAdminExperience)
			admin.PUT("/experience/:id", handlers.UpdateAdminExperience)
			admin.DELETE("/experience/:id", handlers.DeleteAdminExperience)
			admin.POST("/certifications", handlers.CreateAdminCertification)
			admin.PUT("/certifications/:id", handlers.UpdateAdminCertification)
			admin.DELETE("/certifications/:id", handlers.DeleteAdminCertification)
			admin.GET("/messages", handlers.GetAdminMessages)
			admin.GET("/profile", handlers.GetAdminProfile)
			admin.PUT("/profile", handlers.UpdateAdminProfile)
			admin.GET("/skills", handlers.GetAdminSkills)
			admin.POST("/skills", handlers.CreateAdminSkill)
			admin.PUT("/skills/:id", handlers.UpdateAdminSkill)
			admin.DELETE("/skills/:id", handlers.DeleteAdminSkill)
		}
	}

	return router
}

// corsMiddleware mengizinkan frontend Next.js lokal mengakses API dari browser.
// Browser mengirim preflight OPTIONS untuk beberapa request, jadi middleware ini menjawabnya lebih awal
// dengan status 204 agar request asli bisa dilanjutkan oleh browser.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://benny-porto.vercel.app")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		// Request OPTIONS adalah preflight CORS dari browser, bukan endpoint bisnis aplikasi.
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
