package middleware

import (
	"net/http"
	"strings"
	"pamagi/domain"
	"pamagi/internal/tokenutil"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		
		// Pastikan formatnya "Bearer <token>"
		if len(t) == 2 {
			authToken := t[1]
			
			// Validasi token pakai fungsi yang sudah kita buat
			authorized, err := tokenutil.IsAuthorized(authToken, secret)
			if authorized {
				// Ekstrak ID User dari token
				userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
					c.Abort()
					return
				}
				
				// Simpan ID User ke dalam context untuk dipakai di Controller
				c.Set("x-user-id", userID)
				c.Next()
				return
			}
			
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
			c.Abort()
			return
		}
		
		// Jika token tidak ada atau format salah
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Tidak ada akses (Token tidak valid)"})
		c.Abort()
	}
}