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
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid input data format. Please check your request payload."})
		return
	}

	err := wc.WordUsecase.CreateWord(c.Request.Context(), userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to create word. Please try again later."})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Word added successfully"})
}

// GetWords godoc
// @Summary Ambil Daftar Kosakata
// @Description Menampilkan daftar kosakata dengan fitur filter kombo (limit 20)
// @Tags Words
// @Security ApiKeyAuth
// @Produce json
// @Param category_id query string false "Filter by Category ID"
// @Param part_of_speech query string false "Filter by Part of Speech (e.g., NOUN)"
// @Param is_favorite query boolean false "Filter by Favorite status"
// @Param is_bookmarked query boolean false "Filter by Bookmark status"
// @Param start_date query string false "Filter by Start Date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by End Date (YYYY-MM-DD)"
// @Param sort_by query string false "Sorting (newest, oldest, a_z, z_a)"
// @Param page query int false "Nomor Halaman (Default: 1)"
// @Param limit query int false "Jumlah Data per Halaman (Default: 20)"
// @Success 200 {object} dto.WordPaginationResponse
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

// GetWordDetail godoc
// @Summary Ambil Detail Kosakata
// @Description Mengambil satu data kosakata beserta relasi kategori dan contoh kalimatnya
// @Tags Words
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Kosakata"
// @Success 200 {object} dto.WordResponse
// @Failure 404 {object} domain.ErrorResponse
// @Router /words/{id} [get]
func (wc *WordController) GetWordDetail(c *gin.Context) {
	userID := c.GetString("x-user-id")
	wordID := c.Param("id")

	response, err := wc.WordUsecase.GetWordDetail(c.Request.Context(), wordID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Word detail not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ToggleFavorite godoc
// @Summary Ubah Status Favorit
// @Description Menandai atau menghapus tanda favorit pada suatu kata
// @Tags Words
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Kosakata"
// @Success 200 {object} domain.SuccessResponse
// @Router /words/{id}/favorite [patch]
func (wc *WordController) ToggleFavorite(c *gin.Context) {
	userID := c.GetString("x-user-id")
	wordID := c.Param("id")

	err := wc.WordUsecase.ToggleFavorite(c.Request.Context(), wordID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to toggle favorite status"})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Favorite status successfully updated"})
}


// DeleteWord godoc
// @Summary Hapus Kosakata
// @Description Menghapus satu kosakata beserta relasi kategori dan contoh kalimatnya secara permanen
// @Tags Words
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Kosakata"
// @Success 200 {object} domain.SuccessResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /words/{id} [delete]
func (wc *WordController) DeleteWord(c *gin.Context) {
	userID := c.GetString("x-user-id")
	wordID := c.Param("id")

	if err := wc.WordUsecase.DeleteWord(c.Request.Context(), wordID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to delete word. Please ensure the data exists."})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Word deleted successfully"})
}



// GetWordTypes godoc
// @Summary Ambil Daftar Part of Speech
// @Description Menampilkan daftar jenis kata beserta jumlah kata yang dimiliki user
// @Tags Words
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.WordTypeCountResponse
// @Router /words/types [get]
func (wc *WordController) GetWordTypes(c *gin.Context) {
	userID := c.GetString("x-user-id")

	responses, err := wc.WordUsecase.GetWordTypes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to fetch word types"})
		return
	}

	c.JSON(http.StatusOK, responses)
}

// UpdateWord godoc
// @Summary Edit Kosakata
// @Description Memperbarui data kosakata, kategori, dan contoh kalimat
// @Tags Words
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Kosakata"
// @Param request body dto.UpdateWordRequest true "Data Kosakata yang Diupdate"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Router /words/{id} [put]
func (wc *WordController) UpdateWord(c *gin.Context) {
	userID := c.GetString("x-user-id")
	wordID := c.Param("id")
	var request dto.UpdateWordRequest // Harus UpdateWordRequest

if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid input data format. Please check your request payload."})
		return
	}

	if err := wc.WordUsecase.UpdateWord(c.Request.Context(), wordID, userID, &request); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to update word. Ensure the word exists and belongs to you."})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Word updated successfully"})
}