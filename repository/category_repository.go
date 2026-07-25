package repository

import (
	"context"
	"pamagi/domain"
	"pamagi/domain/entity"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(c context.Context, category *entity.Category) error {
	return r.db.WithContext(c).Create(category).Error
}

func (r *categoryRepository) FetchByUserID(c context.Context, userID string) ([]entity.Category, error) {
	var categories []entity.Category
	
	err := r.db.WithContext(c).Table("categories").
		Select("categories.*, COUNT(word_categories.word_id) as count").
		Joins("LEFT JOIN word_categories ON categories.id = word_categories.category_id").
		Where("categories.user_id = ?", userID).
		Group("categories.id").
		Find(&categories).Error
		
	return categories, err
}

// Tambahkan di bagian bawah file category_repository.go
func (r *categoryRepository) CountUncategorized(c context.Context, userID string) (int64, error) {
	var count int64
	// Hitung kata (words) yang TIDAK ADA di tabel word_categories (NULL)
	err := r.db.WithContext(c).Table("words").
		Joins("LEFT JOIN word_categories wc ON words.id = wc.word_id").
		Where("words.user_id = ? AND wc.word_id IS NULL AND words.deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

// Tambahkan di bagian bawah file category_repository.go

func (r *categoryRepository) Update(c context.Context, category *entity.Category) error {
	// Pastikan hanya kategori milik user yang bersangkutan yang di-update
	return r.db.WithContext(c).Model(&entity.Category{}).
		Where("id = ? AND user_id = ?", category.ID, category.UserID).
		Updates(map[string]interface{}{
			"name": category.Name,
			"icon": category.Icon,
		}).Error
}

func (r *categoryRepository) Delete(c context.Context, categoryID string, userID string) error {
	// GORM akan melakukan soft-delete atau hard-delete tergantung konfigurasi modelmu
	// Pastikan menghapus berdasarkan ID dan UserID untuk keamanan
	return r.db.WithContext(c).
		Where("id = ? AND user_id = ?", categoryID, userID).
		Delete(&entity.Category{}).Error
}