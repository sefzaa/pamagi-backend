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
	var examples []entity.WordExample
	for _, exReq := range req.Examples {
		examples = append(examples, entity.WordExample{
			ID: uuid.New().String(), 
			WordID: wordID, 
			RussianSentence: exReq.RussianSentence, 
			TranslatedSentence: exReq.TranslatedSentence,
		})
	}

	word := &entity.Word{
		ID: wordID, 
		UserID: userID, 
		RussianWord: req.RussianWord, 
		Translation: req.Translation, 
		PartOfSpeech: req.PartOfSpeech, 
		Examples: examples,
	}
	return u.wordRepo.Create(c, word, req.CategoryIDs)
}

// Helper untuk mapping Entity ke DTO
func mapToWordResponse(w entity.Word) dto.WordResponse {
	var catRes []dto.CategoryRes
	for _, cat := range w.Categories {
		catRes = append(catRes, dto.CategoryRes{
			ID: cat.ID, 
			Name: cat.Name, 
			Icon: cat.Icon})
	}
	var exRes []dto.WordExampleRes
	for _, ex := range w.Examples {
		exRes = append(exRes, dto.WordExampleRes{
			ID: ex.ID, 
			RussianSentence: ex.RussianSentence, 
			TranslatedSentence: ex.TranslatedSentence})
	}
	return dto.WordResponse{
		ID: w.ID, 
		RussianWord: w.RussianWord, 
		Translation: w.Translation, 
		PartOfSpeech: w.PartOfSpeech,
		IsFavorite: w.IsFavorite, 
		IsBookmarked: w.IsBookmarked, 
		CreatedAt: w.CreatedAt.Format("2006-01-02 15:04:05"),
		Categories: catRes, 
		Examples: exRes,
	}
}

func (u *wordUsecase) GetWords(c context.Context, userID string, filter dto.WordFilterRequest) ([]dto.WordResponse, error) {
	words, err := u.wordRepo.Fetch(c, userID, filter)
	if err != nil { return nil, err }

	var responses []dto.WordResponse
	for _, w := range words {
		responses = append(responses, mapToWordResponse(w))
	}
	if responses == nil { responses = []dto.WordResponse{} }
	return responses, nil
}

func (u *wordUsecase) GetWordDetail(c context.Context, wordID string, userID string) (dto.WordResponse, error) {
	word, err := u.wordRepo.FetchByID(c, wordID, userID)
	if err != nil { return dto.WordResponse{}, err }
	return mapToWordResponse(word), nil
}

func (u *wordUsecase) ToggleFavorite(c context.Context, wordID string, userID string) error {
	return u.wordRepo.ToggleFavorite(c, wordID, userID)
}

func (u *wordUsecase) ToggleBookmark(c context.Context, wordID string, userID string) error {
	return u.wordRepo.ToggleBookmark(c, wordID, userID)
}

func (u *wordUsecase) DeleteWord(c context.Context, wordID string, userID string) error {
	return u.wordRepo.Delete(c, wordID, userID)
}