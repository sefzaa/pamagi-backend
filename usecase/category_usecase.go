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

// Ubah fungsi GetCategories
func (u *categoryUsecase) GetCategories(c context.Context, userID string, forDropdown bool) ([]dto.CategoryResponse, error) {
	categories, err := u.categoryRepo.FetchByUserID(c, userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.CategoryResponse
	for _, cat := range categories {
		responses = append(responses, dto.CategoryResponse{
			ID:    cat.ID,
			Name:  cat.Name,
			Icon:  cat.Icon,
			Count: cat.Count, 
		})
	}

	// JIKA BUKAN UNTUK DROPDOWN, TAMBAHKAN UNCATEGORIZED
	if !forDropdown {
		uncategorizedCount, _ := u.categoryRepo.CountUncategorized(c, userID)
		responses = append(responses, dto.CategoryResponse{
			ID:    "uncategorized", // ID unik untuk FE
			Name:  "Uncategorized",
			Icon:  "folder_off",    // Icon default untuk uncategorized
			Count: uncategorizedCount,
		})
	}

	if responses == nil {
		responses = []dto.CategoryResponse{}
	}
	return responses, nil
}

// Tambahkan di bagian bawah file category_usecase.go

func (u *categoryUsecase) UpdateCategory(c context.Context, categoryID string, userID string, req *dto.CreateCategoryRequest) error {
	category := &entity.Category{
		ID:     categoryID,
		UserID: userID,
		Name:   req.Name,
		Icon:   req.Icon,
	}

	return u.categoryRepo.Update(c, category)
}

func (u *categoryUsecase) DeleteCategory(c context.Context, categoryID string, userID string) error {
	return u.categoryRepo.Delete(c, categoryID, userID)
}