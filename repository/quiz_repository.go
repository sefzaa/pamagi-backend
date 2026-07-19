package repository

import (
	"context"
	"pamagi/domain"
	"pamagi/domain/dto"
	"pamagi/domain/entity"

	"gorm.io/gorm"
)

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) domain.QuizRepository {
	return &quizRepository{db: db}
}

func (r *quizRepository) GenerateQuestions(c context.Context, userID string, filter *dto.GenerateFlashcardRequest, limit int) ([]entity.Word, error) {
	var words []entity.Word
	
	// Mulai query dasar
	query := r.db.WithContext(c).Model(&entity.Word{}).Where("user_id = ?", userID)

	// 1. Terapkan Filter Berdasarkan Pilihan User
	switch filter.FilterType {
	case "category":
		if len(filter.CategoryIDs) > 0 {
			query = query.Joins("JOIN word_categories wc ON wc.word_id = words.id").Where("wc.category_id IN ?", filter.CategoryIDs)
		}
	case "part_of_speech":
		if len(filter.PartsOfSpeech) > 0 {
			query = query.Where("part_of_speech IN ?", filter.PartsOfSpeech)
		}	
	case "favorite":
		query = query.Where("is_favorite = ?", true)
	case "date_range":
		// Format tanggal YYYY-MM-DD
		query = query.Where("DATE(created_at) BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	case "manual":
		query = query.Where("words.id IN ?", filter.WordIDs)
	}

// 2. Hitung jumlah data SEBELUM dilimit
	var count int64
	query.Count(&count)

// 3. Logika Sorting
	switch filter.SortBy {
	case "name_asc":
		query = query.Order("native_word ASC")  // Ubah russian_word menjadi native_word
	case "name_desc":
		query = query.Order("native_word DESC") // Ubah russian_word menjadi native_word
	case "created_asc":
		query = query.Order("created_at ASC")
	case "created_desc":
		query = query.Order("created_at DESC")
	case "updated_desc":
		query = query.Order("updated_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 4. Terapkan Limit dan Tarik Datanya beserta semua relasi Induk-Anak
	// UBAH BARIS INI
	err := query.Preload("Categories").Preload("Targets").Preload("Examples.Targets").Limit(limit).Find(&words).Error	
	return words, err
}


func (r *quizRepository) SaveQuizHistory(c context.Context, session *entity.QuizHistory) error {
	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// 1. Simpan Header (Gunakan Omit agar GORM tidak menyimpan relasi Details secara otomatis)
		if err := tx.Omit("Details").Save(session).Error; err != nil {
			return err
		}

		// 2. HAPUS SEMUA DETAIL LAMA YANG TERKAIT DENGAN SESSION INI
		if err := tx.Exec("DELETE FROM quiz_details WHERE quiz_history_id = ?", session.ID).Error; err != nil {
			return err
		}

		// 3. Masukkan Detail Baru secara manual
		if len(session.Details) > 0 {
			if err := tx.Create(&session.Details).Error; err != nil {
				return err
			}
		}
		return nil
	})
}


func (r *quizRepository) GetHistories(c context.Context, userID string) ([]entity.QuizHistory, error) {
	var sessions []entity.QuizHistory
	err := r.db.WithContext(c).Where("user_id = ?", userID).
		Preload("Details.Word.Targets").          // TAMBAHAN
		Preload("Details.Word.Examples.Targets"). // UBAH BAGIAN INI
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *quizRepository) DeleteHistory(c context.Context, sessionID string, userID string) error {
	// Ubah di bagian Delete(&entity.QuizHistory{})
	return r.db.WithContext(c).Where("id = ? AND user_id = ?", sessionID, userID).Delete(&entity.QuizHistory{}).Error
}