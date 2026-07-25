package domain

import (
	"context"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
)

type CategoryRepository interface {
	Create(c context.Context, category *entity.Category) error
	FetchByUserID(c context.Context, userID string) ([]entity.Category, error)
	CountUncategorized(c context.Context, userID string) (int64, error) // Tambahan
	Update(c context.Context, category *entity.Category) error
	Delete(c context.Context, categoryID string, userID string) error
}



type CategoryUsecase interface {
	// Mengembalikan CategoryResponse agar FE langsung dapat ID-nya
	CreateCategory(c context.Context, userID string, req *dto.CreateCategoryRequest) (dto.CategoryResponse, error)
	GetCategories(c context.Context, userID string, forDropdown bool) ([]dto.CategoryResponse, error) // Tambah param forDropdown
	UpdateCategory(c context.Context, categoryID string, userID string, req *dto.CreateCategoryRequest) error
	DeleteCategory(c context.Context, categoryID string, userID string) error
}

