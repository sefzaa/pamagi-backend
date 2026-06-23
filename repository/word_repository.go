package repository

import (
	"context"
	"pamagi/domain"
	"pamagi/domain/dto"
	"pamagi/domain/entity"

	"gorm.io/gorm"
)

type wordRepository struct {
	db *gorm.DB
}

func NewWordRepository(db *gorm.DB) domain.WordRepository {
	return &wordRepository{db: db}
}

func (r *wordRepository) Create(c context.Context, word *entity.Word, categoryIDs []string) error {
	// Gunakan Transaction agar aman (kalau insert category gagal, insert word juga dibatalkan)
	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// 1. Insert Word beserta Examples-nya (GORM otomatis handle relasi bersarang ini)
		if err := tx.Create(word).Error; err != nil {
			return err
		}

		// 2. Jika ada Kategori yang dipilih, hubungkan ke tabel word_categories
		if len(categoryIDs) > 0 {
			var categories []entity.Category
			tx.Where("id IN ?", categoryIDs).Find(&categories)
			if err := tx.Model(word).Association("Categories").Append(&categories); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *wordRepository) Fetch(c context.Context, userID string, filter dto.WordFilterRequest) ([]entity.Word, error) {
	var words []entity.Word
	
	// Mulai query dengan memanggil relasi (Preload)
	query := r.db.WithContext(c).Preload("Categories").Preload("Examples").Where("user_id = ?", userID)

	// Terapkan Filter Dinamis
	if filter.PartOfSpeech != "" {
		query = query.Where("part_of_speech = ?", filter.PartOfSpeech)
	}
	
	if filter.IsFavorite != nil {
		query = query.Where("is_favorite = ?", *filter.IsFavorite)
	}

	if filter.CategoryID != "" {
		// Join ke tabel pivot jika ingin mencari berdasarkan kategori
		query = query.Joins("JOIN word_categories wc ON wc.word_id = words.id").Where("wc.category_id = ?", filter.CategoryID)
	}

	// Terapkan Sorting
	if filter.SortBy == "oldest" {
		query = query.Order("words.created_at ASC")
	} else {
		query = query.Order("words.created_at DESC") // Default: newest
	}

	err := query.Find(&words).Error
	return words, err
}

func (r *wordRepository) ToggleFavorite(c context.Context, wordID string, userID string) error {
	var word entity.Word
	// Pastikan user hanya bisa mengubah status favorit milik dia sendiri
	if err := r.db.WithContext(c).Where("id = ? AND user_id = ?", wordID, userID).First(&word).Error; err != nil {
		return err
	}

	// Balikkan nilainya (kalau true jadi false, kalau false jadi true)
	word.IsFavorite = !word.IsFavorite
	return r.db.WithContext(c).Save(&word).Error
}