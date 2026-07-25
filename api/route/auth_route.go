package route

import (
	"pamagi/api/controller"
	"pamagi/bootstrap"
	"pamagi/repository"
	"pamagi/usecase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewAuthRouter bertugas merakit fitur Auth dan mendaftarkan endpoint-nya
func NewAuthRouter(env *bootstrap.Env, db *gorm.DB, redis *redis.Client, publicRouter *gin.RouterGroup, protectedRouter *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepository, env, redis)
	authController := &controller.AuthController{
		AuthUsecase: authUsecase,
	}

	// Endpoint publik
	publicRouter.POST("/register", authController.Register)
	publicRouter.POST("/login", authController.Login)
	publicRouter.POST("/refresh", authController.Refresh)

	// Endpoint yang dilindungi satpam (butuh token)
	protectedRouter.POST("/logout", authController.Logout)
	protectedRouter.GET("/users/me", authController.GetProfile)
	protectedRouter.PUT("/users/me", authController.UpdateProfile)
}