package domain

import (
	"context"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
)

type WordRepository interface {
	Create(c context.Context, word *entity.Word, categoryIDs []string) error
	FetchByID(c context.Context, wordID string, userID string) (entity.Word, error) // Baru
	ToggleFavorite(c context.Context, wordID string, userID string) error
	ToggleBookmark(c context.Context, wordID string, userID string) error // Baru
	Delete(c context.Context, wordID string, userID string) error
	Update(c context.Context, word *entity.Word, categoryIDs []string) error

	Fetch(c context.Context, userID string, filter dto.WordFilterRequest) ([]entity.Word, int64, error) // Ubah
	CountByPartOfSpeech(c context.Context, userID string) ([]dto.WordTypeCountResponse, error) // Tambah
}

type WordUsecase interface {
	CreateWord(c context.Context, userID string, req *dto.CreateWordRequest) error
	GetWordDetail(c context.Context, wordID string, userID string) (dto.WordResponse, error) // Baru
	ToggleFavorite(c context.Context, wordID string, userID string) error
	ToggleBookmark(c context.Context, wordID string, userID string) error // Baru
	DeleteWord(c context.Context, wordID string, userID string) error
	UpdateWord(c context.Context, wordID string, userID string, req *dto.CreateWordRequest) error

	GetWords(c context.Context, userID string, filter dto.WordFilterRequest) (dto.WordPaginationResponse, error) // Ubah
	GetWordTypes(c context.Context, userID string) ([]dto.WordTypeCountResponse, error) // Tambah
}


