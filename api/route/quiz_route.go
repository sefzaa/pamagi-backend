package route

import (
	"pamagi/api/controller"
	"pamagi/repository"
	"pamagi/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewQuizRouter mendaftarkan rute untuk fitur Flashcard & Kuis
func NewQuizRouter(db *gorm.DB, protectedRouter *gin.RouterGroup) {
	quizRepository := repository.NewQuizRepository(db)
	quizUsecase := usecase.NewQuizUsecase(quizRepository)
	quizController := &controller.QuizController{
		QuizUsecase: quizUsecase,
	}

	// Semua fitur kuis wajib login
	protectedRouter.POST("/flashcards/generate", quizController.GenerateFlashcards)
	protectedRouter.POST("/flashcards/submit", quizController.SubmitQuiz)
	protectedRouter.GET("/flashcards/history", quizController.GetHistories)
	protectedRouter.DELETE("/flashcards/history/:id", quizController.DeleteHistory)
}