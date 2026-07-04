package controller

import (
	"net/http"
	"pamagi/domain"
	"pamagi/domain/dto"

	"github.com/gin-gonic/gin"
)

type NoteController struct {
	NoteUsecase domain.NoteUsecase
}

// Create godoc
// @Summary Buat Catatan Baru
// @Description Menyimpan catatan baru beserta tag-nya (mendukung format rich text dari FE)
// @Tags Note
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.NoteRequest true "Data Catatan"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Router /notes [post]
func (nc *NoteController) Create(c *gin.Context) {
	userID := c.GetString("x-user-id")
	var req dto.NoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if err := nc.NoteUsecase.CreateNote(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Gagal menyimpan catatan"})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Catatan berhasil disimpan"})
}

// Fetch godoc
// @Summary Ambil Daftar Catatan
// @Description Menampilkan semua catatan milik user (diurutkan dari yang dipin dan terbaru)
// @Param sort query string false "Sort by (title_asc, created_desc, updated_desc)"
// @Tags Note
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.NoteResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /notes [get]
func (nc *NoteController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")
	sortBy := c.Query("sort") // Ambil query misal: /notes?sort=updated_desc

	notes, err := nc.NoteUsecase.GetNotes(c.Request.Context(), userID, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Gagal mengambil catatan"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// Update godoc
// @Summary Edit Catatan
// @Description Memperbarui isi catatan dan tag berdasarkan ID catatan
// @Tags Note
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Catatan"
// @Param request body dto.NoteRequest true "Data Catatan yang Diupdate"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Router /notes/{id} [put]
func (nc *NoteController) Update(c *gin.Context) {
	userID := c.GetString("x-user-id")
	noteID := c.Param("id")
	var req dto.NoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if err := nc.NoteUsecase.UpdateNote(c.Request.Context(), noteID, userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Gagal memperbarui catatan"})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Catatan berhasil diperbarui"})
}

// Delete godoc
// @Summary Hapus Catatan
// @Description Menghapus catatan berdasarkan ID
// @Tags Note
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Catatan"
// @Success 200 {object} domain.SuccessResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /notes/{id} [delete]
func (nc *NoteController) Delete(c *gin.Context) {
	userID := c.GetString("x-user-id")
	noteID := c.Param("id")

	if err := nc.NoteUsecase.DeleteNote(c.Request.Context(), noteID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Gagal menghapus catatan"})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Catatan berhasil dihapus"})
}

// 2. Tambah fungsi GetDetail baru
// GetDetail godoc
// @Summary Ambil Detail Catatan
// @Description Mengambil satu catatan beserta isi lengkapnya
// @Tags Note
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Catatan"
// @Success 200 {object} dto.NoteDetailResponse
// @Failure 404 {object} domain.ErrorResponse
// @Router /notes/{id} [get]
func (nc *NoteController) GetDetail(c *gin.Context) {
	userID := c.GetString("x-user-id")
	noteID := c.Param("id")

	note, err := nc.NoteUsecase.GetNoteDetail(c.Request.Context(), noteID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Catatan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// Tambahkan di bagian bawah file

// TogglePin godoc
// @Summary Toggle Pin Catatan
// @Description Menyematkan atau melepas sematan catatan (hanya mengirim ID)
// @Tags Note
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Catatan"
// @Success 200 {object} domain.SuccessResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /notes/{id}/pin [patch]
func (nc *NoteController) TogglePin(c *gin.Context) {
	userID := c.GetString("x-user-id")
	noteID := c.Param("id")

	if err := nc.NoteUsecase.TogglePin(c.Request.Context(), noteID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Gagal mengubah status pin"})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Status pin berhasil diubah"})
}

// ToggleFavorite godoc
// @Summary Toggle Favorite Catatan
// @Description Menandai atau menghapus tanda favorit pada catatan (hanya mengirim ID)
// @Tags Note
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Catatan"
// @Success 200 {object} domain.SuccessResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /notes/{id}/favorite [patch]
func (nc *NoteController) ToggleFavorite(c *gin.Context) {
	userID := c.GetString("x-user-id")
	noteID := c.Param("id")

	if err := nc.NoteUsecase.ToggleFavorite(c.Request.Context(), noteID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Gagal mengubah status favorit"})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Status favorit berhasil diubah"})
}