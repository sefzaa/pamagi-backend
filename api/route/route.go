package route

import (
	"pamagi/bootstrap"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup meregistrasikan semua route aplikasi
func Setup(env *bootstrap.Env, db *gorm.DB, gin *gin.Engine) {
	// Nanti route spesifik (seperti user_route atau merchant_route)
	// akan kita panggil di dalam sini.
}