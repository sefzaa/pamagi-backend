package entity

import "time"

type User struct {
	ID                 string    `gorm:"type:varchar(36);primaryKey"`
	Name               string    `gorm:"type:varchar(255);not null"`
	Username           string    `gorm:"type:varchar(255);unique;not null"`
	Email              string    `gorm:"type:varchar(255);unique;not null"`
	PasswordHash       string    `gorm:"type:varchar(255);not null"`
	SubscriptionStatus string    `gorm:"type:enum('FREE', 'PREMIUM_LITE', 'PREMIUM_PRO', 'PREMIUM_PLATINUM');default:'FREE'"`
	NoWa               string    `gorm:"type:varchar(20);default:null"`
	NativeLanguage     *string   `gorm:"type:varchar(100);default:null"` // Pengganti Region
	NativeFlagIcon     *string   `gorm:"type:varchar(255);default:null"` // Icon bendera negara asal
	Slogan             string    `gorm:"type:varchar(255);default:'Consistency is key to fluency.'"`
	CreatedAt          time.Time
	UpdatedAt          time.Time

	// Relasi ke bahasa yang dipelajari
	TargetLanguages []UserTargetLanguage `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}