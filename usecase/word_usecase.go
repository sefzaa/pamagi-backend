package usecase

import (
	"context"
	"pamagi/domain"
	"pamagi/domain/dto"
	"pamagi/domain/entity"

	"github.com/google/uuid"
)

type wordUsecase struct {
	wordRepo domain.WordRepository
}

func NewWordUsecase(wordRepo domain.WordRepository) domain.WordUsecase {
	return &wordUsecase{wordRepo: wordRepo}
}

func (u *wordUsecase) CreateWord(c context.Context, userID string, req *dto.CreateWordRequest) error {
	wordID := uuid.New().String()

	// Siapkan data Contoh Kalimat (jika ada)
	var examples []entity.WordExample
	for _, exReq := range req.Examples {
		examples = append(examples, entity.WordExample{
			ID:                 uuid.New().String(),
			WordID:             wordID,
			RussianSentence:    exReq.RussianSentence,
			TranslatedSentence: exReq.TranslatedSentence,
		})
	}

	// Siapkan data Word
	word := &entity.Word{
		ID:           wordID,
		UserID:       userID,
		RussianWord:  req.RussianWord,
		Translation:  req.Translation,
		PartOfSpeech: req.PartOfSpeech,
		Examples:     examples,
	}

	return u.wordRepo.Create(c, word, req.CategoryIDs)
}

func (u *wordUsecase) GetWords(c context.Context, userID string, filter dto.WordFilterRequest) ([]dto.WordResponse, error) {
	words, err := u.wordRepo.Fetch(c, userID, filter)
	if err != nil {
		return nil, err
	}

	// Konversi Entity MySQL menjadi DTO Response yang bersih
	var responses []dto.WordResponse
	for _, w := range words {
		// Mapping Kategori
		var catRes []dto.CategoryRes
		for _, cat := range w.Categories {
			catRes = append(catRes, dto.CategoryRes{
				ID:   cat.ID,
				Name: cat.Name,
			})
		}

		// Mapping Contoh Kalimat
		var exRes []dto.WordExampleRes
		for _, ex := range w.Examples {
			exRes = append(exRes, dto.WordExampleRes{
				ID:                 ex.ID,
				RussianSentence:    ex.RussianSentence,
				TranslatedSentence: ex.TranslatedSentence,
			})
		}

		responses = append(responses, dto.WordResponse{
			ID:           w.ID,
			RussianWord:  w.RussianWord,
			Translation:  w.Translation,
			PartOfSpeech: w.PartOfSpeech,
			IsFavorite:   w.IsFavorite,
			CreatedAt:    w.CreatedAt.Format("2006-01-02 15:04:05"),
			Categories:   catRes,
			Examples:     exRes,
		})
	}

	// Cegah balasan 'null' di JSON jika data kosong
	if responses == nil {
		responses = []dto.WordResponse{}
	}

	return responses, nil
}

func (u *wordUsecase) ToggleFavorite(c context.Context, wordID string, userID string) error {
	return u.wordRepo.ToggleFavorite(c, wordID, userID)
}