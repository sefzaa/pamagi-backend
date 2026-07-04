package domain

import (
	"context"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
)

type WordRepository interface {
	Create(c context.Context, word *entity.Word, categoryIDs []string) error
	Fetch(c context.Context, userID string, filter dto.WordFilterRequest) ([]entity.Word, error)
	FetchByID(c context.Context, wordID string, userID string) (entity.Word, error) // Baru
	ToggleFavorite(c context.Context, wordID string, userID string) error
	ToggleBookmark(c context.Context, wordID string, userID string) error // Baru
	Delete(c context.Context, wordID string, userID string) error
}

type WordUsecase interface {
	CreateWord(c context.Context, userID string, req *dto.CreateWordRequest) error
	GetWords(c context.Context, userID string, filter dto.WordFilterRequest) ([]dto.WordResponse, error)
	GetWordDetail(c context.Context, wordID string, userID string) (dto.WordResponse, error) // Baru
	ToggleFavorite(c context.Context, wordID string, userID string) error
	ToggleBookmark(c context.Context, wordID string, userID string) error // Baru
	DeleteWord(c context.Context, wordID string, userID string) error
}