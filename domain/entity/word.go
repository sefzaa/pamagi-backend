package entity

import "time"

// 1. Tabel Induk: Menyimpan konsep utama (Kata dalam bahasa ibu)
type Word struct {
	ID           string        `gorm:"type:varchar(36);primaryKey"`
	UserID       string        `gorm:"type:varchar(36);not null"`
	NativeWord   string        `gorm:"type:varchar(255);not null"`
	PartOfSpeech string        `gorm:"type:enum('NOUN', 'VERB', 'ADJECTIVE', 'ADVERB', 'PRONOUN', 'PREPOSITION', 'CONJUNCTION', 'INTERJECTION', 'IDIOM');not null"`
	IsFavorite   bool          `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time

	// Relasi
	Categories []Category    `gorm:"many2many:word_categories;constraint:OnDelete:CASCADE;"`
	Targets    []WordTarget  `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
	Examples   []WordExample `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}

// 2. Tabel Anak 1: Menyimpan terjemahan kata per bahasa
type WordTarget struct {
	ID           string `gorm:"type:varchar(36);primaryKey"`
	WordID       string `gorm:"type:varchar(36);not null"`
	LanguageCode string `gorm:"type:varchar(10);not null"`
	TargetWord   string `gorm:"type:varchar(255);not null"`
}

// 3. Tabel Anak 2: Menyimpan kalimat induk (Bahasa ibu)
type WordExample struct {
	ID             string              `gorm:"type:varchar(36);primaryKey"`
	WordID         string              `gorm:"type:varchar(36);not null"`
	NativeSentence string              `gorm:"type:text;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
	
	// Relasi ke terjemahan kalimat
	Targets        []WordExampleTarget `gorm:"foreignKey:WordExampleID;constraint:OnDelete:CASCADE;"`
}

// 4. Tabel Cucu: Menyimpan terjemahan kalimat per bahasa
type WordExampleTarget struct {
	ID             string `gorm:"type:varchar(36);primaryKey"`
	WordExampleID  string `gorm:"type:varchar(36);not null"`
	LanguageCode   string `gorm:"type:varchar(10);not null"`
	TargetSentence string `gorm:"type:text;not null"`
}