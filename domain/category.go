package domain

import (
	"context"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
)

type CategoryRepository interface {
	Create(c context.Context, category *entity.Category) error
	FetchByUserID(c context.Context, userID string) ([]entity.Category, error)
}

type CategoryUsecase interface {
	// Mengembalikan CategoryResponse agar FE langsung dapat ID-nya
	CreateCategory(c context.Context, userID string, req *dto.CreateCategoryRequest) (dto.CategoryResponse, error)
	GetCategories(c context.Context, userID string) ([]dto.CategoryResponse, error)
}