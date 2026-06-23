package controller

import (
	"net/http"
	"pamagi/domain"
	"pamagi/domain/dto"

	"github.com/gin-gonic/gin"
)

type WordController struct {
	WordUsecase domain.WordUsecase
}

// CreateWord godoc
// @Summary Tambah Kosakata Baru
// @Description Menyimpan kosakata beserta contoh kalimat opsional
// @Tags Words
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateWordRequest true "Data Kosakata"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Router /words [post]
func (wc *WordController) CreateWord(c *gin.Context) {
	userID := c.GetString("x-user-id")
	var request dto.CreateWordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err := wc.WordUsecase.CreateWord(c.Request.Context(), userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Kosakata berhasil ditambahkan"})
}

// GetWords godoc
// @Summary Ambil Daftar Kosakata
// @Description Menampilkan daftar kosakata dengan fitur filter
// @Tags Words
// @Security ApiKeyAuth
// @Produce json
// @Param category_id query string false "Filter by Category ID"
// @Param part_of_speech query string false "Filter by Part of Speech (e.g., NOUN)"
// @Param is_favorite query boolean false "Filter by Favorite status"
// @Param sort_by query string false "Sorting (newest, oldest)"
// @Success 200 {array} dto.WordResponse
// @Router /words [get]
func (wc *WordController) GetWords(c *gin.Context) {
	userID := c.GetString("x-user-id")
	var filter dto.WordFilterRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	responses, err := wc.WordUsecase.GetWords(c.Request.Context(), userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}

// ToggleFavorite godoc
// @Summary Ubah Status Favorit
// @Description Menandai atau menghapus tanda favorit pada suatu kata
// @Tags Words
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Kosakata"
// @Success 200 {object} domain.SuccessResponse
// @Router /words/{id}/favorite [put]
func (wc *WordController) ToggleFavorite(c *gin.Context) {
	userID := c.GetString("x-user-id")
	wordID := c.Param("id")

	err := wc.WordUsecase.ToggleFavorite(c.Request.Context(), wordID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Gagal merubah status favorit"})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Status favorit diperbarui"})
}