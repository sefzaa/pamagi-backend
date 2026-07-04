package entity

import "time"

type Note struct {
	ID         string    `gorm:"type:varchar(36);primaryKey"`
	UserID     string    `gorm:"type:varchar(36);not null"`
	Title      string    `gorm:"type:varchar(255);not null"`
	Content    string    `gorm:"type:text;not null"`
	Color      string    `gorm:"type:varchar(10);default:'#FFFFFF'"`
	IsPinned   bool      `gorm:"default:false"`
	IsFavorite bool      `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Relasi many-to-many ke tabel ntags melalui tabel pivot note_tags
	Tags []Ntag `gorm:"many2many:note_tags;joinForeignKey:NoteID;joinReferences:NtagID"`
}

type Ntag struct {
	ID        string `gorm:"type:varchar(36);primaryKey"`
	Name      string `gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time
}