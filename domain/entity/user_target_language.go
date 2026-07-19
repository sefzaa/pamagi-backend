package entity

import "time"

type UserTargetLanguage struct {
	ID           string `gorm:"type:varchar(36);primaryKey"`
	UserID       string `gorm:"type:varchar(36);not null"`
	LanguageCode string `gorm:"type:varchar(10);not null"`
	LanguageName string `gorm:"type:varchar(100);not null"`
	FlagIcon     string `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time
}