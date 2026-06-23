package route

import (
	"pamagi/api/controller"
	"pamagi/repository"
	"pamagi/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewWordRouter bertugas merakit fitur Word dan mendaftarkan endpoint-nya
func NewWordRouter(db *gorm.DB, protectedRouter *gin.RouterGroup) {
	wordRepository := repository.NewWordRepository(db)
	wordUsecase := usecase.NewWordUsecase(wordRepository)
	wordController := &controller.WordController{
		WordUsecase: wordUsecase,
	}

	// Semua fitur word harus login dulu, jadi kita pakai protectedRouter
	protectedRouter.POST("/words", wordController.CreateWord)
	protectedRouter.GET("/words", wordController.GetWords)
	protectedRouter.PUT("/words/:id/favorite", wordController.ToggleFavorite)
}