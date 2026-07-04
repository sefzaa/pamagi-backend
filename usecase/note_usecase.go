package usecase

import (
	"context"
	"pamagi/domain"
	"pamagi/domain/dto"
	"pamagi/domain/entity"

	"github.com/google/uuid"
)

type noteUsecase struct {
	noteRepo domain.NoteRepository
}

func NewNoteUsecase(noteRepo domain.NoteRepository) domain.NoteUsecase {
	return &noteUsecase{noteRepo: noteRepo}
}

func (u *noteUsecase) CreateNote(c context.Context, userID string, req *dto.NoteRequest) error {
	// 1. Proses tags terlebih dahulu
	tags, err := u.noteRepo.FindOrCreateTags(c, req.Tags)
	if err != nil {
		return err
	}

	// 2. Susun entity Note
	note := &entity.Note{
		ID:         uuid.New().String(),
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		Color:      req.Color,
		IsPinned:   req.IsPinned,
		IsFavorite: req.IsFavorite,
		Tags:       tags,
	}

	// 3. Simpan
	return u.noteRepo.Create(c, note)
}

func (u *noteUsecase) DeleteNote(c context.Context, noteID string, userID string) error {
	return u.noteRepo.Delete(c, noteID, userID)
}

func (u *noteUsecase) UpdateNote(c context.Context, noteID string, userID string, req *dto.NoteRequest) error {
	// 1. Proses tags yang diinputkan user saat edit
	tags, err := u.noteRepo.FindOrCreateTags(c, req.Tags)
	if err != nil {
		return err
	}

	// 2. Susun entity Note
	note := &entity.Note{
		ID:         noteID,
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		Color:      req.Color,
		IsPinned:   req.IsPinned,
		IsFavorite: req.IsFavorite,
		Tags:       tags,
	}

	// 3. Simpan pembaruan
	return u.noteRepo.Update(c, note)
}

// 1. Ubah GetNotes
func (u *noteUsecase) GetNotes(c context.Context, userID string, sortBy string) ([]dto.NoteListResponse, error) {
	notes, err := u.noteRepo.FetchByUserID(c, userID, sortBy)
	if err != nil {
		return nil, err
	}

	var responses []dto.NoteListResponse
	for _, n := range notes {
		var tagNames []string
		for _, t := range n.Tags {
			tagNames = append(tagNames, t.Name)
		}

		// Buat preview (potong maksimal 50 karakter). Kita pakai rune agar karakter Cyrillic bahasa Rusia aman dipotong.
		runes := []rune(n.Content)
		preview := n.Content
		if len(runes) > 50 {
			preview = string(runes[:50]) + "..."
		}

		responses = append(responses, dto.NoteListResponse{
			ID:         n.ID,
			Title:      n.Title,
			Preview:    preview, // Kirim potongannya saja
			Color:      n.Color,
			IsPinned:   n.IsPinned,
			IsFavorite: n.IsFavorite,
			Tags:       tagNames,
			UpdatedAt:  n.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if responses == nil {
		responses = []dto.NoteListResponse{}
	}
	return responses, nil
}

// 2. Tambah GetNoteDetail
func (u *noteUsecase) GetNoteDetail(c context.Context, noteID string, userID string) (dto.NoteDetailResponse, error) {
	n, err := u.noteRepo.FetchByID(c, noteID, userID)
	if err != nil {
		return dto.NoteDetailResponse{}, err
	}

	var tagNames []string
	for _, t := range n.Tags {
		tagNames = append(tagNames, t.Name)
	}

	return dto.NoteDetailResponse{
		ID:         n.ID,
		Title:      n.Title,
		Content:    n.Content, // Kirim teks utuh
		Color:      n.Color,
		IsPinned:   n.IsPinned,
		IsFavorite: n.IsFavorite,
		Tags:       tagNames,
		UpdatedAt:  n.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// Tambahkan di bagian bawah file
func (u *noteUsecase) TogglePin(c context.Context, noteID string, userID string) error {
	return u.noteRepo.TogglePin(c, noteID, userID)
}

func (u *noteUsecase) ToggleFavorite(c context.Context, noteID string, userID string) error {
	return u.noteRepo.ToggleFavorite(c, noteID, userID)
}