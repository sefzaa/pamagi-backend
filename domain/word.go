package domain

import (
	"context"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
)

// Repository: Mengurus query ke database MySQL
type WordRepository interface {
	Create(c context.Context, word *entity.Word, categoryIDs []string) error
	Fetch(c context.Context, userID string, filter dto.WordFilterRequest) ([]entity.Word, error)
	ToggleFavorite(c context.Context, wordID string, userID string) error
	// Nanti kita bisa tambah Update dan Delete di sini
}

// Usecase: Mengurus logika bisnis
type WordUsecase interface {
	CreateWord(c context.Context, userID string, req *dto.CreateWordRequest) error
	GetWords(c context.Context, userID string, filter dto.WordFilterRequest) ([]dto.WordResponse, error)
	ToggleFavorite(c context.Context, wordID string, userID string) error
}