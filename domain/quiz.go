package domain

import (
	"context"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
)

type QuizRepository interface {
	// Untuk menarik kata-kata berdasarkan filter yang dinamis
	GenerateQuestions(c context.Context, userID string, filter *dto.GenerateFlashcardRequest, limit int) ([]entity.Word, error)
	
	// Untuk menyimpan histori kuis
	SaveQuizSession(c context.Context, session *entity.QuizHistory) error
	
	// Untuk fitur Premium
	GetHistories(c context.Context, userID string) ([]entity.QuizHistory, error)
	DeleteHistory(c context.Context, sessionID string, userID string) error

}

type QuizUsecase interface {
	GenerateFlashcards(c context.Context, userID string, req *dto.GenerateFlashcardRequest) ([]dto.WordResponse, error)
	SubmitQuiz(c context.Context, userID string, req *dto.SubmitQuizRequest) error
	
	GetQuizHistories(c context.Context, userID string) ([]dto.QuizSessionResponse, error)
	DeleteQuizHistory(c context.Context, sessionID string, userID string) error
}