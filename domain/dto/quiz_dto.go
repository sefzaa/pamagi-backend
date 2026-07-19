package dto

// 1. Request untuk FE meminta soal Flashcard
type GenerateFlashcardRequest struct {
	FilterType   string   `json:"filter_type" binding:"required"` // Opsi: "category", "part_of_speech", "favorite", "date_range", "manual", "all"
	
	// Parameter opsional (tergantung FilterType)
	CategoryIDs   []string `json:"category_ids"`
	PartsOfSpeech []string `json:"parts_of_speech"`
	StartDate    string   `json:"start_date"` // Format: YYYY-MM-DD
	EndDate      string   `json:"end_date"`
	WordIDs      []string `json:"word_ids"`   // Jika pilih manual satu per satu
	TotalQuestions int `json:"total_questions" binding:"required,min=1,max=100"`
	
	SortBy       string   `json:"sort_by"`    // Opsi: "name_asc", "name_desc", "created_asc", "created_desc", "updated_desc"
}

// 2. Request untuk FE mengirim laporan setelah selesai kuis
type SubmitQuizRequest struct {
	ID               string             `json:"id"`
	TotalWords       int                `json:"total_words" binding:"required"`
	CorrectAnswers   int                `json:"correct_answers"`
	IncorrectAnswers int                `json:"incorrect_answers"`
	Score            float64            `json:"score"`
	Status           string             `json:"status" binding:"required,oneof=IN_PROGRESS COMPLETED"` // Tambahkan ini
	Details          []SubmitQuizDetail `json:"details"`
}

type SubmitQuizDetail struct {
	WordID    string `json:"word_id" binding:"required"`
	IsCorrect *bool  `json:"is_correct"`
}

// 3. Response untuk histori kuis
type QuizHistoryResponse struct {
	ID               string             `json:"id"`
	TotalWords       int                `json:"total_words"`
	CorrectAnswers   int                `json:"correct_answers"`
	IncorrectAnswers int                `json:"incorrect_answers"`
	Score            float64            `json:"score"`
	Status           string               `json:"status"`
	CreatedAt        string             `json:"created_at"`
	Details          []QuizDetailResponse `json:"details,omitempty"` // omitempty agar null untuk user Free
}

type QuizDetailResponse struct {
	WordID             string           `json:"word_id"`
	TargetLanguageCode string           `json:"target_language_code"` // TAMBAHAN
	TargetWord         string           `json:"target_word"`          // PENGGANTI RussianWord
	NativeWord         string           `json:"native_word"`          // PENGGANTI Translation
	PartOfSpeech       string           `json:"part_of_speech"`
	IsCorrect          *bool            `json:"is_correct"`
	Examples           []WordExampleRes `json:"examples"`
}


