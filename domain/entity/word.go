package entity

import "time"

type Word struct {
	ID                 string        `gorm:"type:varchar(36);primaryKey"`
	UserID             string        `gorm:"type:varchar(36);not null"`
	TargetLanguageCode string        `gorm:"type:varchar(10);not null"` // Tambahan untuk membedakan bahasa
	TargetWord         string        `gorm:"type:varchar(255);not null"` // Pengganti RussianWord
	NativeWord         string        `gorm:"type:varchar(255);not null"` // Pengganti Translation
	PartOfSpeech       string        `gorm:"type:enum('NOUN', 'VERB', 'ADJECTIVE', 'ADVERB', 'PRONOUN', 'PREPOSITION', 'CONJUNCTION', 'INTERJECTION', 'IDIOM');not null"`
	IsFavorite         bool          `gorm:"default:false"`
	IsBookmarked       bool          `gorm:"default:false"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time    

	Categories   []Category    `gorm:"many2many:word_categories;constraint:OnDelete:CASCADE;"`
	Examples     []WordExample `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}

type WordExample struct {
	ID             string     `gorm:"type:varchar(36);primaryKey"`
	WordID         string     `gorm:"type:varchar(36);not null"`
	TargetSentence string     `gorm:"type:text;not null"` // Pengganti RussianSentence
	NativeSentence string     `gorm:"type:text;not null"` // Pengganti TranslatedSentence
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}