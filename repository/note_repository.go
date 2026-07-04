package repository

import (
	"context"
	"pamagi/domain"
	"pamagi/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) domain.NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) Create(c context.Context, note *entity.Note) error {
	return r.db.WithContext(c).Create(note).Error
}


func (r *noteRepository) FindOrCreateTags(c context.Context, tagNames []string) ([]entity.Ntag, error) {
	var tags []entity.Ntag
	for _, name := range tagNames {
		var tag entity.Ntag
		// Cek apakah tag sudah ada, jika belum, buat baru
		err := r.db.WithContext(c).Where("name = ?", name).FirstOrCreate(&tag, entity.Ntag{
			ID:   uuid.New().String(),
			Name: name,
		}).Error
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *noteRepository) Delete(c context.Context, noteID string, userID string) error {
	return r.db.WithContext(c).Where("id = ? AND user_id = ?", noteID, userID).Delete(&entity.Note{}).Error
}

func (r *noteRepository) Update(c context.Context, note *entity.Note) error {
	// 1. Update data utama menggunakan map agar nilai boolean (false) tetap ter-update
	err := r.db.WithContext(c).Model(note).Where("id = ? AND user_id = ?", note.ID, note.UserID).Updates(map[string]interface{}{
		"title":       note.Title,
		"content":     note.Content,
		"color":       note.Color,
		"is_pinned":   note.IsPinned,
		"is_favorite": note.IsFavorite,
	}).Error
	if err != nil {
		return err
	}

	// 2. Ganti semua relasi Tags lama dengan Tags yang baru di tabel pivot
	return r.db.WithContext(c).Model(note).Association("Tags").Replace(note.Tags)
}

// 1. Ubah FetchByUserID
func (r *noteRepository) FetchByUserID(c context.Context, userID string, sortBy string) ([]entity.Note, error) {
	var notes []entity.Note
	query := r.db.WithContext(c).Preload("Tags").Where("user_id = ?", userID)

	// Logika Pinned + Sorting Dinamis
	switch sortBy {
	case "title_asc":
		query = query.Order("is_pinned DESC, title ASC")
	case "title_desc":
		query = query.Order("is_pinned DESC, title DESC")
	case "created_asc":
		query = query.Order("is_pinned DESC, created_at ASC")
	case "created_desc":
		query = query.Order("is_pinned DESC, created_at DESC")
	case "updated_desc":
		query = query.Order("is_pinned DESC, updated_at DESC")
	default:
		// Default: yang di-pin di atas, lalu urutkan berdasarkan yang paling baru diedit
		query = query.Order("is_pinned DESC, updated_at DESC") 
	}

	err := query.Find(&notes).Error
	return notes, err
}

// 2. Tambah FetchByID (Simpan di bawah FetchByUserID)
func (r *noteRepository) FetchByID(c context.Context, noteID string, userID string) (entity.Note, error) {
	var note entity.Note
	err := r.db.WithContext(c).Preload("Tags").Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error
	return note, err
}

// Tambahkan di bagian bawah file
func (r *noteRepository) TogglePin(c context.Context, noteID string, userID string) error {
	return r.db.WithContext(c).Model(&entity.Note{}).
		Where("id = ? AND user_id = ?", noteID, userID).
		UpdateColumn("is_pinned", gorm.Expr("NOT is_pinned")).Error
}

func (r *noteRepository) ToggleFavorite(c context.Context, noteID string, userID string) error {
	return r.db.WithContext(c).Model(&entity.Note{}).
		Where("id = ? AND user_id = ?", noteID, userID).
		UpdateColumn("is_favorite", gorm.Expr("NOT is_favorite")).Error
}