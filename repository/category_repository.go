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
	// Ambil kategori yang hanya dimiliki oleh user yang sedang login
	err := r.db.WithContext(c).Where("user_id = ?", userID).Find(&categories).Error
	return categories, err
}