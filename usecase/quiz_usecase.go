package usecase

import (
	"context"
	"math/rand" // Tambahkan ini
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
    isPremium := false 

    limit := req.TotalQuestions
    if limit > 20 && !isPremium {
        limit = 20 
    }

    words, err := u.quizRepo.GenerateQuestions(c, userID, req, limit)
    if err != nil {
        return nil, err
    }

    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    r.Shuffle(len(words), func(i, j int) {
        words[i], words[j] = words[j], words[i]
    })

	var responses []dto.WordResponse
	for _, w := range words {
		var catRes []dto.CategoryRes
		for _, cat := range w.Categories {
			catRes = append(catRes, dto.CategoryRes{ID: cat.ID, Name: cat.Name})
		}
		
		var targetRes []dto.WordTargetRes
		for _, t := range w.Targets {
			targetRes = append(targetRes, dto.WordTargetRes{
				ID:           t.ID,
				LanguageCode: t.LanguageCode,
				TargetWord:   t.TargetWord,
			})
		}

		var exRes []dto.WordExampleRes
		for _, ex := range w.Examples {
			var extRes []dto.ExampleTargetRes
			for _, ext := range ex.Targets {
				extRes = append(extRes, dto.ExampleTargetRes{
					ID:           ext.ID,
					LanguageCode: ext.LanguageCode,
					Sentence:     ext.TargetSentence,
				})
			}
			exRes = append(exRes, dto.WordExampleRes{
				ID:              ex.ID,
				NativeSentence:  ex.NativeSentence,
				TargetSentences: extRes,
			})
		}

		responses = append(responses, dto.WordResponse{
			ID:           w.ID, 
			NativeWord:   w.NativeWord,
			PartOfSpeech: w.PartOfSpeech, 
			IsFavorite:   w.IsFavorite,
			CreatedAt:    w.CreatedAt.Format("2006-01-02 15:04:05"),
			Categories:   catRes,
			Targets:      targetRes,
			Examples:     exRes,
		})
	}

	if responses == nil {
		responses = []dto.WordResponse{}
	}
	return responses, nil
}

func (u *quizUsecase) SubmitQuiz(c context.Context, userID string, req *dto.SubmitQuizRequest) error {
	sessionID := req.ID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	// 1. Buat Header
	session := &entity.QuizHistory{
		ID:               sessionID,
		UserID:           userID,
		TotalQuestions:   req.TotalWords,
		CorrectAnswers:   req.CorrectAnswers,
		IncorrectAnswers: req.IncorrectAnswers,
		Score:            req.Score,
		Status:           req.Status,
		CreatedAt:        time.Now(),
	}

	// 2. Siapkan detail DENGAN UUID BARU SETIAP KALI DIPANGGIL
	var details []entity.QuizDetail
	for _, d := range req.Details {
		details = append(details, entity.QuizDetail{
			ID:            uuid.New().String(), // INI KUNCI UTAMA: Setiap loop buat UUID baru
			QuizHistoryID: sessionID,
			WordID:        d.WordID,
			IsCorrect:     d.IsCorrect,
		})
	}
	session.Details = details 

	return u.quizRepo.SaveQuizHistory(c, session)
}

func (u *quizUsecase) GetQuizHistories(c context.Context, userID string) ([]dto.QuizHistoryResponse, error) {
	sessions, err := u.quizRepo.GetHistories(c, userID)
	if err != nil {
		return nil, err
	}

	isPremium := true 

	var responses []dto.QuizHistoryResponse
	for _, s := range sessions {
		var detailRes []dto.QuizDetailResponse
		
		if isPremium {
			for _, d := range s.Details {
				
				var targetRes []dto.WordTargetRes
				for _, t := range d.Word.Targets {
					targetRes = append(targetRes, dto.WordTargetRes{
						ID:           t.ID,
						LanguageCode: t.LanguageCode,
						TargetWord:   t.TargetWord,
					})
				}

				var exRes []dto.WordExampleRes
				for _, ex := range d.Word.Examples {
					var extRes []dto.ExampleTargetRes
					for _, ext := range ex.Targets {
						extRes = append(extRes, dto.ExampleTargetRes{
							ID:           ext.ID,
							LanguageCode: ext.LanguageCode,
							Sentence:     ext.TargetSentence,
						})
					}
					exRes = append(exRes, dto.WordExampleRes{
						ID:              ex.ID,
						NativeSentence:  ex.NativeSentence,
						TargetSentences: extRes,
					})
				}

				detailRes = append(detailRes, dto.QuizDetailResponse{
					WordID:       d.WordID,
					NativeWord:   d.Word.NativeWord,
					PartOfSpeech: d.Word.PartOfSpeech,
					IsCorrect:    d.IsCorrect,
					Targets:      targetRes,
					Examples:     exRes,
				})
			}
		}

		responses = append(responses, dto.QuizHistoryResponse{
			ID:               s.ID,
			TotalWords:       s.TotalQuestions,
			CorrectAnswers:   s.CorrectAnswers,
			IncorrectAnswers: s.IncorrectAnswers,
			Score:            s.Score,
			Status:           s.Status,
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