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
		query = query.Joins("JOIN word_categories wc ON wc.word_id = words.id").Where("wc.category_id = ?", filter.CategoryID)
	case "part_of_speech":
		query = query.Where("part_of_speech = ?", filter.PartOfSpeech)
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

	// 3. Logika Sorting & Randomizer
	if int(count) > limit {
		// Jika data melebihi limit (misal > 20), paksa urutan jadi ACAK
		query = query.Order("RAND()")
	} else {
		// Jika data kurang dari limit, ikuti kemauan sorting user
		switch filter.SortBy {
		case "name_asc":
			query = query.Order("russian_word ASC")
		case "name_desc":
			query = query.Order("russian_word DESC")
		case "created_asc":
			query = query.Order("created_at ASC")
		case "created_desc":
			query = query.Order("created_at DESC")
		case "updated_desc":
			query = query.Order("updated_at DESC")
		default:
			query = query.Order("created_at DESC") // Default terbaru
		}
	}

	// 4. Terapkan Limit dan Tarik Datanya beserta relasinya
	err := query.Preload("Categories").Preload("Examples").Limit(limit).Find(&words).Error
	return words, err
}


func (r *quizRepository) SaveQuizSession(c context.Context, session *entity.QuizHistory) error {
	return r.db.WithContext(c).Create(session).Error
}


func (r *quizRepository) GetHistories(c context.Context, userID string) ([]entity.QuizHistory, error) {
	var sessions []entity.QuizHistory // Ubah di sini
	err := r.db.WithContext(c).Where("user_id = ?", userID).
		Preload("Details.Word"). 
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *quizRepository) DeleteHistory(c context.Context, sessionID string, userID string) error {
	// Ubah di bagian Delete(&entity.QuizHistory{})
	return r.db.WithContext(c).Where("id = ? AND user_id = ?", sessionID, userID).Delete(&entity.QuizHistory{}).Error
}