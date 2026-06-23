package dto

// Format request saat FE kirim nama kategori baru
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// Format balasan untuk dikirim kembali ke FE
type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}