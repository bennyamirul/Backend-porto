package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AdminAuthMiddleware melindungi route admin dengan memvalidasi JWT dari cookie admin_token.
// JWT adalah token bertanda tangan: server tidak perlu menyimpan sesi di database, cukup mengecek
// signature memakai ADMIN_JWT_SECRET dan mengecek expiry token. Parameter masuknya adalah context Gin
// dari setiap request, dan hasilnya adalah melanjutkan handler berikutnya atau menghentikan request
// dengan response 401 kalau cookie tidak ada, token rusak, signature salah, atau token kedaluwarsa.
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("admin_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Pastikan token memakai algoritma HMAC seperti HS256; ini mencegah token dengan algoritma lain diterima.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return []byte(os.Getenv("ADMIN_JWT_SECRET")), nil
		})

		if err != nil || !token.Valid || claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}
