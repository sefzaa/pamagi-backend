package controller

import (
	"net/http"
	"pamagi/domain"
	"pamagi/domain/dto"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	CategoryUsecase domain.CategoryUsecase
}

// CreateCategory godoc
// @Summary Tambah Kategori Baru
// @Description Menyimpan kategori baru milik user (misal: "Kata Kerja A1")
// @Tags Categories
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Data Kategori"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} domain.ErrorResponse
// @Router /categories [post]
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	userID := c.GetString("x-user-id")
	var request dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	response, err := cc.CategoryUsecase.CreateCategory(c.Request.Context(), userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCategories godoc
// @Summary Ambil Daftar Kategori
// @Description Menampilkan semua daftar kategori.
// @Tags Categories
// @Security ApiKeyAuth
// @Produce json
// @Param for_dropdown query boolean false "Set true untuk menyembunyikan Uncategorized (digunakan saat input kata)"
// @Success 200 {array} dto.CategoryResponse
// @Router /categories [get]
func (cc *CategoryController) GetCategories(c *gin.Context) {
	userID := c.GetString("x-user-id")
	
	// Cek apakah FE mengirim ?for_dropdown=true
	forDropdown := c.Query("for_dropdown") == "true"

	responses, err := cc.CategoryUsecase.GetCategories(c.Request.Context(), userID, forDropdown)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}


// UpdateCategory godoc
// @Summary Edit Kategori
// @Description Memperbarui nama atau ikon kategori milik user
// @Tags Categories
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Kategori"
// @Param request body dto.CreateCategoryRequest true "Data Kategori Baru"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /categories/{id} [put]
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	userID := c.GetString("x-user-id")
	categoryID := c.Param("id")
	var request dto.CreateCategoryRequest // Gunakan DTO yang sama atau buat UpdateCategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid input data format. Please check your request."})
		return
	}

	// Catatan: Pastikan fungsi UpdateCategory sudah dibuat di CategoryUsecase nanti
	if err := cc.CategoryUsecase.UpdateCategory(c.Request.Context(), categoryID, userID, &request); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to update category. Ensure the category exists and belongs to you."})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Category updated successfully"})
}

// DeleteCategory godoc
// @Summary Hapus Kategori
// @Description Menghapus kategori secara permanen
// @Tags Categories
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Kategori"
// @Success 200 {object} domain.SuccessResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /categories/{id} [delete]
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	userID := c.GetString("x-user-id")
	categoryID := c.Param("id")

	// Catatan: Pastikan fungsi DeleteCategory sudah dibuat di CategoryUsecase nanti
	if err := cc.CategoryUsecase.DeleteCategory(c.Request.Context(), categoryID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to delete category. Please try again later."})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Category deleted successfully"})
}