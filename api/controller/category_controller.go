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
// @Description Menampilkan semua daftar kategori milik user yang sedang login
// @Tags Categories
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.CategoryResponse
// @Router /categories [get]
func (cc *CategoryController) GetCategories(c *gin.Context) {
	userID := c.GetString("x-user-id")

	responses, err := cc.CategoryUsecase.GetCategories(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}