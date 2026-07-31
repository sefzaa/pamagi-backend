package route

import (
	"pamagi/api/controller"
	"pamagi/bootstrap"
	"pamagi/internal/mailer" // IMPORT BARU
	"pamagi/repository"
	"pamagi/usecase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewAuthRouter(env *bootstrap.Env, db *gorm.DB, redis *redis.Client, publicRouter *gin.RouterGroup, protectedRouter *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(db)
	
	// Inisialisasi Email Service
	//emailService := mailer.NewBrevoService(env)
	emailService := mailer.NewGmailService(env)
	
	// Inject Email Service ke Usecase
	authUsecase := usecase.NewAuthUsecase(authRepository, env, redis, emailService)
	
	authController := &controller.AuthController{
		AuthUsecase: authUsecase,
	}

	// Endpoint publik
	publicRouter.Static("/assets", "./assets")
	publicRouter.POST("/register", authController.Register)
	publicRouter.POST("/login", authController.Login)
	publicRouter.POST("/refresh", authController.Refresh)
	publicRouter.POST("/forgot-password", authController.ForgotPassword) // ENDPOINT BARU
	publicRouter.POST("/verify-otp", authController.VerifyOTP)         // TAMBAHAN BARU
	publicRouter.POST("/reset-password", authController.ResetPassword)

	// Endpoint yang dilindungi satpam
	protectedRouter.POST("/logout", authController.Logout)
	protectedRouter.GET("/users/me", authController.GetProfile)
	protectedRouter.PUT("/users/me", authController.UpdateProfile)
}