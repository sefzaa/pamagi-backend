package entity

import "time"

type Category struct {
	ID        string `gorm:"type:varchar(36);primaryKey"`
	UserID    string `gorm:"type:varchar(36);not null"`
	Name      string `gorm:"type:varchar(255);not null"`
	Icon      string `gorm:"type:varchar(255);default:'folder'"` 
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}