package repository

import (
	"context"
	"strings"
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

// Cari fungsi Create dan ganti dengan kode ini:
func (r *wordRepository) Create(c context.Context, word *entity.Word, categoryIDs []string) error {
	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Simpan Induk (GORM otomatis akan menyimpan Targets dan Examples)
		if err := tx.Create(word).Error; err != nil {
			return err
		}
		
		// Hubungkan dengan kategori
		if len(categoryIDs) > 0 {
			var categories []entity.Category
			if err := tx.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
				return err
			}
			if err := tx.Model(word).Association("Categories").Append(&categories); err != nil {
				return err
			}
		}
		return nil
	})
}

// 1. Ubah fungsi Fetch
func (r *wordRepository) Fetch(c context.Context, userID string, filter dto.WordFilterRequest) ([]entity.Word, int64, error) {
	var words []entity.Word
	var total int64 
	
	query := r.db.WithContext(c).Model(&entity.Word{}).Where("words.user_id = ?", userID)

	// --- FILTER BAHASA (JOIN ke tabel anak word_targets) ---
	if filter.TargetLanguageCode != "" {
		query = query.Joins("JOIN word_targets wt ON wt.word_id = words.id").Where("wt.language_code = ?", filter.TargetLanguageCode)
	}

	// 1. Filter Part of Speech
	if filter.PartOfSpeech != "" { 
		posList := strings.Split(filter.PartOfSpeech, ",")
		query = query.Where("words.part_of_speech IN ?", posList) 
	}

	if filter.IsFavorite != nil { query = query.Where("words.is_favorite = ?", *filter.IsFavorite) }
	// (Filter IsBookmarked sudah dibuang)
	
	// 2. Filter Category ID
	if filter.CategoryID != "" {
		catList := strings.Split(filter.CategoryID, ",")
		hasUncategorized := false
		var validIDs []string

		for _, id := range catList {
			if id == "uncategorized" { hasUncategorized = true } else { validIDs = append(validIDs, id) }
		}

		if hasUncategorized && len(validIDs) > 0 {
			query = query.Joins("LEFT JOIN word_categories wc ON wc.word_id = words.id").Where("wc.category_id IN ? OR wc.word_id IS NULL", validIDs)
		} else if hasUncategorized {
			query = query.Joins("LEFT JOIN word_categories wc ON wc.word_id = words.id").Where("wc.word_id IS NULL")
		} else if len(validIDs) > 0 {
			query = query.Joins("JOIN word_categories wc ON wc.word_id = words.id").Where("wc.category_id IN ?", validIDs)
		}
	}

	if filter.StartDate != "" && filter.EndDate != "" {
		query = query.Where("DATE(words.created_at) BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	}

	// Group By agar tidak duplikat gara-gara JOIN
	query = query.Group("words.id")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch filter.SortBy {
	case "oldest": query = query.Order("words.created_at ASC")
	case "a_z": query = query.Order("words.native_word ASC")
	case "z_a": query = query.Order("words.native_word DESC")
	default: query = query.Order("words.created_at DESC")
	}

	offset := (filter.Page - 1) * filter.Limit
	// Pastikan kita me-load seluruh cucu relasinya
	err := query.Preload("Categories").Preload("Targets").Preload("Examples.Targets").Limit(filter.Limit).Offset(offset).Find(&words).Error
	
	return words, total, err 
}

// 2. Tambahkan fungsi baru untuk menghitung Part of Speech
func (r *wordRepository) CountByPartOfSpeech(c context.Context, userID string) ([]dto.WordTypeCountResponse, error) {
	var counts []dto.WordTypeCountResponse
	err := r.db.WithContext(c).Model(&entity.Word{}).
		Select("part_of_speech, count(id) as count").
		Where("user_id = ? AND deleted_at IS NULL", userID). // Jangan hitung yang sudah dihapus
		Group("part_of_speech").
		Scan(&counts).Error
	return counts, err
}

func (r *wordRepository) FetchByID(c context.Context, wordID string, userID string) (entity.Word, error) {
	var word entity.Word
	err := r.db.WithContext(c).Preload("Categories").Preload("Targets").Preload("Examples.Targets").Where("id = ? AND user_id = ?", wordID, userID).First(&word).Error
	return word, err
}

func (r *wordRepository) ToggleFavorite(c context.Context, wordID string, userID string) error {
	return r.db.WithContext(c).Model(&entity.Word{}).Where("id = ? AND user_id = ?", wordID, userID).
		UpdateColumn("is_favorite", gorm.Expr("NOT is_favorite")).Error
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
		// 1. Cek existensi
		var existing entity.Word
		if err := tx.Where("id = ? AND user_id = ?", word.ID, word.UserID).First(&existing).Error; err != nil {
			return err
		}

		// 2. Update Induk Utama
		if err := tx.Model(&existing).Updates(map[string]interface{}{
			"native_word":    word.NativeWord,
			"part_of_speech": word.PartOfSpeech,
		}).Error; err != nil {
			return err
		}

		// 3. Update Kategori Pivot
		var categories []entity.Category
		if len(categoryIDs) > 0 { tx.Where("id IN ?", categoryIDs).Find(&categories) }
		if err := tx.Model(&existing).Association("Categories").Replace(&categories); err != nil {
			return err
		}

		// 4. Bersihkan Anak & Cucu Lama
		tx.Where("word_id = ?", word.ID).Delete(&entity.WordTarget{})
		tx.Where("word_id = ?", word.ID).Delete(&entity.WordExample{})
		// (word_example_targets otomatis terhapus karena Cascade di database)

		// 5. Masukkan Anak & Cucu Baru
		if len(word.Targets) > 0 {
			if err := tx.Create(&word.Targets).Error; err != nil { return err }
		}
		if len(word.Examples) > 0 {
			if err := tx.Create(&word.Examples).Error; err != nil { return err }
		}

		return nil
	})
}



