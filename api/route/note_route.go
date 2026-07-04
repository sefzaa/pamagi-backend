package route

import (
	"pamagi/api/controller"
	"pamagi/repository"
	"pamagi/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewNoteRouter mendaftarkan rute untuk fitur Catatan (Notes)
func NewNoteRouter(db *gorm.DB, protectedRouter *gin.RouterGroup) {
	noteRepository := repository.NewNoteRepository(db)
	noteUsecase := usecase.NewNoteUsecase(noteRepository)
	noteController := &controller.NoteController{
		NoteUsecase: noteUsecase,
	}

	// Semua fitur note wajib login, jadi kita pakai protectedRouter
	noteGroup := protectedRouter.Group("/notes")
	{
		noteGroup.POST("", noteController.Create)
		noteGroup.GET("", noteController.Fetch)
		noteGroup.PUT("/:id", noteController.Update)
		noteGroup.GET("/:id", noteController.GetDetail)
		noteGroup.DELETE("/:id", noteController.Delete)

		noteGroup.PATCH("/:id/pin", noteController.TogglePin)
		noteGroup.PATCH("/:id/favorite", noteController.ToggleFavorite)
	}
}