package usecase

import (
	"context"
	"pamagi/domain"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
	"github.com/google/uuid"
)

type categoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(categoryRepo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo}
}

func (u *categoryUsecase) CreateCategory(c context.Context, userID string, req *dto.CreateCategoryRequest) (dto.CategoryResponse, error) {
	category := &entity.Category{
		ID:     uuid.New().String(),
		UserID: userID,
		Name:   req.Name,
		Icon:   req.Icon, // Simpan Icon
	}

	err := u.categoryRepo.Create(c, category)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	return dto.CategoryResponse{
		ID: category.ID, 
		Name: category.Name, 
		Icon: category.Icon}, nil
}

func (u *categoryUsecase) GetCategories(c context.Context, userID string) ([]dto.CategoryResponse, error) {
	categories, err := u.categoryRepo.FetchByUserID(c, userID)
	if err != nil {
		return nil, err
	}
	var responses []dto.CategoryResponse
	for _, cat := range categories {
		responses = append(responses, dto.CategoryResponse{
			ID: cat.ID, 
			Name: cat.Name, 
			Icon: cat.Icon})
	}
	if responses == nil { responses = []dto.CategoryResponse{} }
	return responses, nil
}