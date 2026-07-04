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
	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(word).Error; err != nil {
			return err
		}
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
	query := r.db.WithContext(c).Preload("Categories").Preload("Examples").Where("user_id = ?", userID)

	if filter.PartOfSpeech != "" {
		query = query.Where("part_of_speech = ?", filter.PartOfSpeech)
	}
	if filter.IsFavorite != nil {
		query = query.Where("is_favorite = ?", *filter.IsFavorite)
	}
	if filter.IsBookmarked != nil {
		query = query.Where("is_bookmarked = ?", *filter.IsBookmarked)
	}
	if filter.CategoryID != "" {
		query = query.Joins("JOIN word_categories wc ON wc.word_id = words.id").Where("wc.category_id = ?", filter.CategoryID)
	}
	if filter.StartDate != "" && filter.EndDate != "" {
		query = query.Where("DATE(words.created_at) BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	}

	switch filter.SortBy {
	case "oldest":
		query = query.Order("words.created_at ASC")
	case "a_z":
		query = query.Order("words.russian_word ASC") // Database UTF-8 otomatis urutkan Cyrillic/Kanji
	case "z_a":
		query = query.Order("words.russian_word DESC")
	default:
		query = query.Order("words.created_at DESC") // Default: newest
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.Limit(filter.Limit).Offset(offset).Find(&words).Error
	return words, err
}

func (r *wordRepository) FetchByID(c context.Context, wordID string, userID string) (entity.Word, error) {
	var word entity.Word
	err := r.db.WithContext(c).Preload("Categories").Preload("Examples").Where("id = ? AND user_id = ?", wordID, userID).First(&word).Error
	return word, err
}

func (r *wordRepository) ToggleFavorite(c context.Context, wordID string, userID string) error {
	return r.db.WithContext(c).Model(&entity.Word{}).Where("id = ? AND user_id = ?", wordID, userID).
		UpdateColumn("is_favorite", gorm.Expr("NOT is_favorite")).Error
}

func (r *wordRepository) ToggleBookmark(c context.Context, wordID string, userID string) error {
	return r.db.WithContext(c).Model(&entity.Word{}).Where("id = ? AND user_id = ?", wordID, userID).
		UpdateColumn("is_bookmarked", gorm.Expr("NOT is_bookmarked")).Error
}

func (r *wordRepository) Delete(c context.Context, wordID string, userID string) error {
	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var word entity.Word
		// Pastikan kata ada dan memang milik user yang sedang login
		if err := tx.Where("id = ? AND user_id = ?", wordID, userID).First(&word).Error; err != nil {
			return err
		}

		// 1. Hapus relasi di tabel pivot (word_categories)
		if err := tx.Model(&word).Association("Categories").Clear(); err != nil {
			return err
		}

		// 2. Hapus anak tabel (contoh kalimat)
		if err := tx.Where("word_id = ?", wordID).Delete(&entity.WordExample{}).Error; err != nil {
			return err
		}

		// 3. Hapus kata utamanya
		return tx.Delete(&word).Error
	})
}

func (r *wordRepository) Update(c context.Context, word *entity.Word, categoryIDs []string) error {
	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// 1. Pastikan kata tersebut ada dan memang milik user yang sedang login
		var existing entity.Word
		if err := tx.Where("id = ? AND user_id = ?", word.ID, word.UserID).First(&existing).Error; err != nil {
			return err
		}

		// 2. Update data utama kata
		if err := tx.Model(&existing).Updates(map[string]interface{}{
			"russian_word":   word.RussianWord,
			"translation":    word.Translation,
			"part_of_speech": word.PartOfSpeech,
		}).Error; err != nil {
			return err
		}

		// 3. Replace relasi Kategori (tabel pivot word_categories otomatis diurus GORM)
		var categories []entity.Category
		if len(categoryIDs) > 0 {
			tx.Where("id IN ?", categoryIDs).Find(&categories)
		}
		if err := tx.Model(&existing).Association("Categories").Replace(&categories); err != nil {
			return err
		}

		// 4. Update Contoh Kalimat (Hapus yang lama, simpan yang baru)
		if err := tx.Where("word_id = ?", word.ID).Delete(&entity.WordExample{}).Error; err != nil {
			return err
		}
		if len(word.Examples) > 0 {
			// Pastikan setiap contoh kalimat baru di-binding ke WordID yang sedang diedit
			for i := range word.Examples {
				word.Examples[i].WordID = word.ID
			}
			if err := tx.Create(&word.Examples).Error; err != nil {
				return err
			}
		}

		return nil
	})
}