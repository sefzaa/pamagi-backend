package route

import (
	"pamagi/api/middleware"
	"pamagi/bootstrap"
	"time" // Tambahkan import ini
	"github.com/gin-contrib/cors" // Tambahkan import ini

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, db *gorm.DB, redis *redis.Client, ginEngine *gin.Engine) {
	//TAMBAHKAN MIDDLEWARE CORS DI SINI
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://pamagi.mydm.cloud"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 1. Siapkan Grup Rute Publik
	publicRouter := ginEngine.Group("")

	// 2. Siapkan Grup Rute Private (Dilindungi Middleware JWT)
	protectedRouter := ginEngine.Group("")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	

	// 3. Panggil dan daftarkan semua Router Fitur di sini
	NewAuthRouter(env, db, redis, publicRouter, protectedRouter)
	NewWordRouter(db, protectedRouter)
	NewCategoryRouter(db, protectedRouter)
	NewQuizRouter(db, protectedRouter)
	NewNoteRouter(db, protectedRouter)
	
	
	
	// Nanti kalau ada fitur baru, tinggal tambah di sini:
	// NewQuizRouter(db, protectedRouter)
	// NewProfileRouter(db, protectedRouter)
}