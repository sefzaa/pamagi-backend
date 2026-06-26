package usecase

import (
	"context"
	"time"
	"pamagi/domain"
	"pamagi/domain/dto"
	"pamagi/domain/entity"

	"github.com/google/uuid"
)

type quizUsecase struct {
	quizRepo domain.QuizRepository
}

func NewQuizUsecase(quizRepo domain.QuizRepository) domain.QuizUsecase {
	return &quizUsecase{quizRepo: quizRepo}
}

func (u *quizUsecase) GenerateFlashcards(c context.Context, userID string, req *dto.GenerateFlashcardRequest) ([]dto.WordResponse, error) {
	// TODO: Nanti kita bisa query ke tabel users untuk cek apakah dia PREMIUM.
	// Sementara kita hardcode default Free (20 kata). Jika premium, ubah jadi 100.
	isPremium := false 
	limit := 20
	if isPremium {
		limit = 100
	}

	words, err := u.quizRepo.GenerateQuestions(c, userID, req, limit)
	if err != nil {
		return nil, err
	}

	// Mapping Entity ke DTO WordResponse (Bisa pakai struktur response Word yang sudah ada)
	var responses []dto.WordResponse
	for _, w := range words {
		var catRes []dto.CategoryRes
		for _, cat := range w.Categories {
			catRes = append(catRes, dto.CategoryRes{ID: cat.ID, Name: cat.Name})
		}
		var exRes []dto.WordExampleRes
		for _, ex := range w.Examples {
			exRes = append(exRes, dto.WordExampleRes{
				ID: ex.ID, RussianSentence: ex.RussianSentence, TranslatedSentence: ex.TranslatedSentence,
			})
		}
		responses = append(responses, dto.WordResponse{
			ID: w.ID, RussianWord: w.RussianWord, Translation: w.Translation,
			PartOfSpeech: w.PartOfSpeech, IsFavorite: w.IsFavorite,
			CreatedAt: w.CreatedAt.Format("2006-01-02 15:04:05"),
			Categories: catRes, Examples: exRes,
		})
	}

	if responses == nil {
		responses = []dto.WordResponse{}
	}
	return responses, nil
}

func (u *quizUsecase) SubmitQuiz(c context.Context, userID string, req *dto.SubmitQuizRequest) error {
	sessionID := uuid.New().String()

	var details []entity.QuizDetail
	for _, d := range req.Details {
		details = append(details, entity.QuizDetail{
			ID:            uuid.New().String(),
			QuizHistoryID: sessionID,
			WordID:        d.WordID,
			IsCorrect:     d.IsCorrect,
		})
	}

	session := &entity.QuizHistory{
		ID:               sessionID,
		UserID:           userID,
		TotalQuestions:   req.TotalWords,
		CorrectAnswers:   req.CorrectAnswers,
		IncorrectAnswers: req.IncorrectAnswers,
		Score:            req.Score,
		Status:           req.Status, // Ambil status dari DTO
		Details:          details,
		CreatedAt:        time.Now(),
	}

	return u.quizRepo.SaveQuizHistory(c, session)
}

func (u *quizUsecase) GetQuizHistories(c context.Context, userID string) ([]dto.QuizHistoryResponse, error) {
	sessions, err := u.quizRepo.GetHistories(c, userID)
	if err != nil {
		return nil, err
	}

	// TODO: Cek status Premium user
	isPremium := false 

	var responses []dto.QuizHistoryResponse
	for _, s := range sessions {
		var detailRes []dto.QuizDetailResponse
		
		// Jika Premium, masukkan data salah/benarnya. Jika Free, biarkan kosong.
		if isPremium {
			for _, d := range s.Details {
				detailRes = append(detailRes, dto.QuizDetailResponse{
					WordID:      d.WordID,
					RussianWord: d.Word.RussianWord,
					Translation: d.Word.Translation,
					IsCorrect:   d.IsCorrect,
				})
			}
		}

		responses = append(responses, dto.QuizHistoryResponse{
			ID:               s.ID,
			TotalWords:       s.TotalQuestions,
			CorrectAnswers:   s.CorrectAnswers,
			IncorrectAnswers: s.IncorrectAnswers,
			Score:            s.Score,
			CreatedAt:        s.CreatedAt.Format("2006-01-02 15:04:05"),
			Details:          detailRes,
		})
	}

	if responses == nil {
		responses = []dto.QuizHistoryResponse{}
	}
	return responses, nil
}

func (u *quizUsecase) DeleteQuizHistory(c context.Context, sessionID string, userID string) error {
	return u.quizRepo.DeleteHistory(c, sessionID, userID)
}