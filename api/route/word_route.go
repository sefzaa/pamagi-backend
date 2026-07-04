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

	protectedRouter.GET("/words/types", wordController.GetWordTypes)

	
	protectedRouter.GET("/words/:id", wordController.GetWordDetail) // Detail lengkap
	protectedRouter.DELETE("/words/:id", wordController.DeleteWord) // Hapus kata
	protectedRouter.PUT("/words/:id", wordController.UpdateWord) // Tambahan Edit Kosakata

	// Ubah PUT menjadi PATCH di sini
	protectedRouter.PATCH("/words/:id/favorite", wordController.ToggleFavorite)
	protectedRouter.PATCH("/words/:id/bookmark", wordController.ToggleBookmark)
}
