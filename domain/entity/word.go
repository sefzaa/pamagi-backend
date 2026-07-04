package entity

import "time"

// Representasi tabel words
type Word struct {
	ID           string        `gorm:"type:varchar(36);primaryKey"`
	UserID       string        `gorm:"type:varchar(36);not null"`
	RussianWord  string        `gorm:"type:varchar(255);not null"`
	Translation  string        `gorm:"type:varchar(255);not null"`
	PartOfSpeech string        `gorm:"type:enum('NOUN', 'VERB', 'ADJECTIVE', 'ADVERB', 'PRONOUN', 'PREPOSITION', 'CONJUNCTION', 'INTERJECTION', 'IDIOM');not null"`
	IsFavorite   bool          `gorm:"default:false"`
	IsBookmarked bool          `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time    // Gunakan pointer untuk soft delete

	// Relasi ke tabel lain
	Categories   []Category    `gorm:"many2many:word_categories;constraint:OnDelete:CASCADE;"`
	Examples     []WordExample `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}



// Representasi tabel word_examples
type WordExample struct {
	ID                 string     `gorm:"type:varchar(36);primaryKey"`
	WordID             string     `gorm:"type:varchar(36);not null"`
	RussianSentence    string     `gorm:"type:text;not null"`
	TranslatedSentence string     `gorm:"type:text;not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
}