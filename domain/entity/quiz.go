package entity

import "time"

type QuizHistory struct {
	ID               string `gorm:"type:varchar(36);primaryKey"`
	UserID           string `gorm:"type:varchar(36);not null"`
	TotalQuestions   int    `gorm:"not null"` // Sesuai kolom SQL total_questions
	CorrectAnswers   int    `gorm:"default:0"`
	IncorrectAnswers int    `gorm:"default:0"`
	Score            float64 `gorm:"default:0"`
	Status           string `gorm:"type:enum('IN_PROGRESS', 'COMPLETED');default:'IN_PROGRESS'"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	Details []QuizDetail `gorm:"foreignKey:QuizHistoryID;constraint:OnDelete:CASCADE;"`
}

type QuizDetail struct {
	ID            string `gorm:"type:varchar(36);primaryKey"`
	QuizHistoryID string `gorm:"type:varchar(36);not null"` // Ubah foreign key-nya
	WordID        string `gorm:"type:varchar(36);not null"`
	IsCorrect 	 *bool 	 `gorm:"column:is_correct"`

	Word Word `gorm:"foreignKey:WordID"`
}