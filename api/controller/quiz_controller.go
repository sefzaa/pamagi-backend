package controller

import (
	"net/http"
	"pamagi/domain"
	"pamagi/domain/dto"

	"github.com/gin-gonic/gin"
)

type QuizController struct {
	QuizUsecase domain.QuizUsecase
}

// GenerateFlashcards godoc
// @Summary Dapatkan Soal Flashcard
// @Description Mengambil daftar kata berdasarkan filter (maksimal 20 untuk Free, diacak jika lebih)
// @Tags Flashcard
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.GenerateFlashcardRequest true "Filter Flashcard"
// @Success 200 {array} dto.WordResponse
// @Failure 400 {object} domain.ErrorResponse
// @Router /flashcards/generate [post]
func (qc *QuizController) GenerateFlashcards(c *gin.Context) {
	userID := c.GetString("x-user-id")
	var request dto.GenerateFlashcardRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	responses, err := qc.QuizUsecase.GenerateFlashcards(c.Request.Context(), userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}

// SubmitQuiz godoc
// @Summary Kirim Hasil Kuis
// @Description Menyimpan skor akhir dan detail jawaban user setelah kuis selesai
// @Tags Flashcard
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.SubmitQuizRequest true "Data Skor Kuis"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Router /flashcards/submit [post]
func (qc *QuizController) SubmitQuiz(c *gin.Context) {
	userID := c.GetString("x-user-id")
	var request dto.SubmitQuizRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err := qc.QuizUsecase.SubmitQuiz(c.Request.Context(), userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Quiz results saved successfully"})
}

// GetHistories godoc
// @Summary Ambil Histori Kuis
// @Description Menampilkan riwayat kuis user (Detail jawaban hanya muncul jika Premium)
// @Tags Flashcard
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.QuizHistoryResponse
// @Router /flashcards/history [get]
func (qc *QuizController) GetHistories(c *gin.Context) {
	userID := c.GetString("x-user-id")

	responses, err := qc.QuizUsecase.GetQuizHistories(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}

// DeleteHistory godoc
// @Summary Hapus Histori Kuis
// @Description Menghapus satu riwayat kuis berdasarkan ID sesi
// @Tags Flashcard
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Session"
// @Success 200 {object} domain.SuccessResponse
// @Router /flashcards/history/{id} [delete]
func (qc *QuizController) DeleteHistory(c *gin.Context) {
	userID := c.GetString("x-user-id")
	sessionID := c.Param("id")

	err := qc.QuizUsecase.DeleteQuizHistory(c.Request.Context(), sessionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to delete quiz history"})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Quiz history deleted successfully"})
}