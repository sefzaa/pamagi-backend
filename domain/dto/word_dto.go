package dto

// ==========================================
// REQUEST: FORMAT JSON DARI FRONTEND
// ==========================================

type CreateWordRequest struct {
	NativeWord   string           `json:"native_word"`    // Hapus binding:"required"
	PartOfSpeech string           `json:"part_of_speech"` // Hapus binding:"required"
	CategoryIDs  []string         `json:"category_ids"`
	Targets      []WordTargetReq  `json:"targets" binding:"required,max=2"` // Array-nya wajib ada, tapi isinya boleh kosong
	Examples     []WordExampleReq `json:"examples" binding:"max=3"`
}

type WordTargetReq struct {
	LanguageCode string `json:"language_code" binding:"required"` // Language code tetap wajib agar tahu bendera apa yang ditampilkan
	TargetWord   string `json:"target_word"`                      // Hapus binding:"required"
}

type UpdateWordRequest struct {
	NativeWord   string           `json:"native_word"`    // Hapus binding:"required"
	PartOfSpeech string           `json:"part_of_speech"` // Hapus binding:"required"
	CategoryIDs  []string         `json:"category_ids"`
	Targets      []WordTargetReq  `json:"targets" binding:"required,max=2"`
	Examples     []WordExampleReq `json:"examples" binding:"max=3"`
}

type WordExampleReq struct {
	NativeSentence string             `json:"native_sentence" binding:"required"`
	TargetSentences []ExampleTargetReq `json:"target_sentences" binding:"required"`
}

type ExampleTargetReq struct {
	LanguageCode string `json:"language_code" binding:"required"`
	Sentence     string `json:"sentence" binding:"required"`
}


type WordFilterRequest struct {
	TargetLanguageCode string `form:"target_language_code"`
	CategoryID         string `form:"category_id"`
	PartOfSpeech       string `form:"part_of_speech"`
	IsFavorite         *bool  `form:"is_favorite"`
	StartDate          string `form:"start_date"`
	EndDate            string `form:"end_date"`
	SortBy             string `form:"sort_by"`
	Page               int    `form:"page,default=1"`
	Limit              int    `form:"limit,default=20"`
}

// ==========================================
// RESPONSE: FORMAT JSON KE FRONTEND
// ==========================================

type WordResponse struct {
	ID           string           `json:"id"`
	NativeWord   string           `json:"native_word"`
	PartOfSpeech string           `json:"part_of_speech"`
	IsFavorite   bool             `json:"is_favorite"`
	CreatedAt    string           `json:"created_at"`
	Categories   []CategoryRes    `json:"categories"`
	Targets      []WordTargetRes  `json:"targets"`
	Examples     []WordExampleRes `json:"examples"`
}

type CategoryRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type WordTargetRes struct {
	ID           string `json:"id"`
	LanguageCode string `json:"language_code"`
	TargetWord   string `json:"target_word"`
}

type WordExampleRes struct {
	ID              string             `json:"id"`
	NativeSentence  string             `json:"native_sentence"`
	TargetSentences []ExampleTargetRes `json:"target_sentences"`
}

type ExampleTargetRes struct {
	ID           string `json:"id"`
	LanguageCode string `json:"language_code"`
	Sentence     string `json:"sentence"`
}

type PaginationMeta struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
}

type WordPaginationResponse struct {
	Data []WordResponse `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type WordTypeCountResponse struct {
	PartOfSpeech string `json:"part_of_speech"`
	Count        int64  `json:"count"`
}