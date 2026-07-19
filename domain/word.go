package domain

import (
	"context"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
)

type WordRepository interface {
	Create(c context.Context, word *entity.Word, categoryIDs []string) error
	Fetch(c context.Context, userID string, filter dto.WordFilterRequest) ([]entity.Word, int64, error)
	CountByPartOfSpeech(c context.Context, userID string) ([]dto.WordTypeCountResponse, error)
	FetchByID(c context.Context, wordID string, userID string) (entity.Word, error)
	ToggleFavorite(c context.Context, wordID string, userID string) error
	Delete(c context.Context, wordID string, userID string) error
    // PERBAIKAN DI SINI: Gunakan Update, bukan UpdateWord
	Update(c context.Context, word *entity.Word, categoryIDs []string) error 
}

type WordUsecase interface {
	CreateWord(c context.Context, userID string, req *dto.CreateWordRequest) error
	GetWordDetail(c context.Context, wordID string, userID string) (dto.WordResponse, error) 
	ToggleFavorite(c context.Context, wordID string, userID string) error
	DeleteWord(c context.Context, wordID string, userID string) error
	UpdateWord(c context.Context, wordID string, userID string, req *dto.UpdateWordRequest) error

	GetWords(c context.Context, userID string, filter dto.WordFilterRequest) (dto.WordPaginationResponse, error) 
	GetWordTypes(c context.Context, userID string) ([]dto.WordTypeCountResponse, error) 
}