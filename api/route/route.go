package route

import (
	"pamagi/api/controller"
	"pamagi/bootstrap"
	"pamagi/repository"
	"pamagi/usecase"
	"pamagi/api/middleware"


	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"


)

func Setup(env *bootstrap.Env, db *gorm.DB, redis *redis.Client, gin *gin.Engine) {
	// Inisialisasi Auth Layer dengan menyuntikkan redis client
	authRepository := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepository, env, redis)
	authController := &controller.AuthController{
		AuthUsecase: authUsecase,
	}

	publicRouter := gin.Group("")
	{
		publicRouter.POST("/register", authController.Register)
		publicRouter.POST("/login", authController.Login)
	}


	// 2. Rute Private (Wajib bawa token)
	protectedRouter := gin.Group("")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret)) 
	{
		protectedRouter.POST("/logout", authController.Logout)
	}
}