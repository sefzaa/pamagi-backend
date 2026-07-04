package domain

import (
	"context"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
)

type NoteRepository interface {
	Create(c context.Context, note *entity.Note) error
	FetchByUserID(c context.Context, userID string, sortBy string) ([]entity.Note, error) // Ubah
	FindOrCreateTags(c context.Context, tagNames []string) ([]entity.Ntag, error)
	Delete(c context.Context, noteID string, userID string) error
	Update(c context.Context, note *entity.Note) error

	
	FetchByID(c context.Context, noteID string, userID string) (entity.Note, error)       // Tambah Baru

	TogglePin(c context.Context, noteID string, userID string) error
	ToggleFavorite(c context.Context, noteID string, userID string) error
}

type NoteUsecase interface {
	CreateNote(c context.Context, userID string, req *dto.NoteRequest) error
	GetNotes(c context.Context, userID string, sortBy string) ([]dto.NoteListResponse, error) // Ubah
	DeleteNote(c context.Context, noteID string, userID string) error
	UpdateNote(c context.Context, noteID string, userID string, req *dto.NoteRequest) error // Tambahan Edit

	
	GetNoteDetail(c context.Context, noteID string, userID string) (dto.NoteDetailResponse, error) // Tambah Baru

	TogglePin(c context.Context, noteID string, userID string) error
	ToggleFavorite(c context.Context, noteID string, userID string) error
}