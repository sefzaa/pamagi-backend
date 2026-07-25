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

	// Default ke "NONE" jika kosong
	partOfSpeech := req.PartOfSpeech
	if partOfSpeech == "" {
		partOfSpeech = "NONE"
	}

	// 1. Looping Anak (Target Bahasa)
	var targets []entity.WordTarget
	for _, t := range req.Targets {
		targets = append(targets, entity.WordTarget{
			ID:           uuid.New().String(),
			WordID:       wordID,
			LanguageCode: t.LanguageCode,
			TargetWord:   t.TargetWord,
		})
	}

	// 2. Looping Induk Contoh Kalimat & Cucu (Terjemahan Kalimat)
	var examples []entity.WordExample
	for _, exReq := range req.Examples {
		exID := uuid.New().String()
		
		var exTargets []entity.WordExampleTarget
		for _, extReq := range exReq.TargetSentences {
			exTargets = append(exTargets, entity.WordExampleTarget{
				ID:             uuid.New().String(),
				WordExampleID:  exID,
				LanguageCode:   extReq.LanguageCode,
				TargetSentence: extReq.Sentence,
			})
		}
		
		examples = append(examples, entity.WordExample{
			ID:             exID,
			WordID:         wordID,
			NativeSentence: exReq.NativeSentence,
			Targets:        exTargets,
		})
	}

	// 3. Bungkus menjadi 1 Konsep Utama (Induk)
word := &entity.Word{
		ID:           wordID,
		UserID:       userID,
		NativeWord:   req.NativeWord,
		PartOfSpeech: partOfSpeech, // Gunakan variabel yang sudah di-cek
		Targets:      targets,
		Examples:     examples,
	}

	return u.wordRepo.Create(c, word, req.CategoryIDs)
}

// Helper untuk mapping Entity ke DTO
func mapToWordResponse(w entity.Word) dto.WordResponse {
	var catRes []dto.CategoryRes
	for _, cat := range w.Categories {
		catRes = append(catRes, dto.CategoryRes{ID: cat.ID, Name: cat.Name, Icon: cat.Icon})
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

	return dto.WordResponse{
		ID:           w.ID,
		NativeWord:   w.NativeWord,
		PartOfSpeech: w.PartOfSpeech,
		IsFavorite:   w.IsFavorite,
		CreatedAt:    w.CreatedAt.Format("2006-01-02 15:04:05"),
		Categories:   catRes,
		Targets:      targetRes,
		Examples:     exRes,
	}
}

// 1. Ubah fungsi GetWords
func (u *wordUsecase) GetWords(c context.Context, userID string, filter dto.WordFilterRequest) (dto.WordPaginationResponse, error) {
	words, total, err := u.wordRepo.Fetch(c, userID, filter)
	if err != nil {
		return dto.WordPaginationResponse{}, err
	}

	var responses []dto.WordResponse
	for _, w := range words {
		responses = append(responses, mapToWordResponse(w))
	}
	if responses == nil {
		responses = []dto.WordResponse{}
	}

	// Kalkulasi total halaman
	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))

	return dto.WordPaginationResponse{
		Data: responses,
		Meta: dto.PaginationMeta{
			TotalItems:  total,
			TotalPages:  totalPages,
			CurrentPage: filter.Page,
			Limit:       filter.Limit,
		},
	}, nil
}

// 2. Tambahkan fungsi GetWordTypes
func (u *wordUsecase) GetWordTypes(c context.Context, userID string) ([]dto.WordTypeCountResponse, error) {
	return u.wordRepo.CountByPartOfSpeech(c, userID)
}

func (u *wordUsecase) GetWordDetail(c context.Context, wordID string, userID string) (dto.WordResponse, error) {
	word, err := u.wordRepo.FetchByID(c, wordID, userID)
	if err != nil { return dto.WordResponse{}, err }
	return mapToWordResponse(word), nil
}

func (u *wordUsecase) ToggleFavorite(c context.Context, wordID string, userID string) error {
	return u.wordRepo.ToggleFavorite(c, wordID, userID)
}


func (u *wordUsecase) DeleteWord(c context.Context, wordID string, userID string) error {
	return u.wordRepo.Delete(c, wordID, userID)
}


func (u *wordUsecase) UpdateWord(c context.Context, wordID string, userID string, req *dto.UpdateWordRequest) error {
	// Default ke "NONE" jika kosong
	partOfSpeech := req.PartOfSpeech
	if partOfSpeech == "" {
		partOfSpeech = "NONE"
	}

	var targets []entity.WordTarget
	for _, t := range req.Targets {
		targets = append(targets, entity.WordTarget{
			ID:           uuid.New().String(),
			WordID:       wordID,
			LanguageCode: t.LanguageCode,
			TargetWord:   t.TargetWord,
		})
	}

	var examples []entity.WordExample
	for _, exReq := range req.Examples {
		exID := uuid.New().String()
		var exTargets []entity.WordExampleTarget
		for _, extReq := range exReq.TargetSentences {
			exTargets = append(exTargets, entity.WordExampleTarget{
				ID:             uuid.New().String(),
				WordExampleID:  exID,
				LanguageCode:   extReq.LanguageCode,
				TargetSentence: extReq.Sentence,
			})
		}
		examples = append(examples, entity.WordExample{
			ID:             exID,
			WordID:         wordID,
			NativeSentence: exReq.NativeSentence,
			Targets:        exTargets,
		})
	}
	word := &entity.Word{
		ID:           wordID,
		UserID:       userID,
		NativeWord:   req.NativeWord,
		PartOfSpeech: partOfSpeech, // Gunakan variabel yang sudah di-cek
		Targets:      targets,
		Examples:     examples,
	}

	return u.wordRepo.Update(c, word, req.CategoryIDs)
}