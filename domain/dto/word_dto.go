package dto

// 1. Request untuk Create Word
type CreateWordRequest struct {
	NativeWord   string              `json:"native_word" binding:"required"`
	PartOfSpeech string              `json:"part_of_speech" binding:"required"`
	CategoryIDs  []string            `json:"category_ids"`
	Targets      []WordTargetRequest `json:"targets" binding:"required,min=1,max=2"` // Array bahasa asing yang diinput (maks 2)
}

type WordTargetRequest struct {
	LanguageCode string           `json:"language_code" binding:"required"`
	TargetWord   string           `json:"target_word" binding:"required"`
	Examples     []ExampleRequest `json:"examples" binding:"max=3"` // Maksimal 3 contoh kalimat per bahasa target
}

// Sub-request untuk contoh kalimat
type ExampleRequest struct {
	TargetSentence string `json:"target_sentence" binding:"required"`
	NativeSentence string `json:"native_sentence" binding:"required"`
}

// 2. Request untuk Filter GET /words (Pake Query Params)
type WordFilterRequest struct {
	TargetLanguageCode string `form:"target_language_code"` // TAMBAHAN: Filter berdasar bahasa (contoh: "en", "ru")
	CategoryID         string `form:"category_id"`
	PartOfSpeech       string `form:"part_of_speech"`
	IsFavorite         *bool  `form:"is_favorite"`
	IsBookmarked       *bool  `form:"is_bookmarked"`
	StartDate          string `form:"start_date"`
	EndDate            string `form:"end_date"`
	SortBy             string `form:"sort_by"`
	Page               int    `form:"page,default=1"`
	Limit              int    `form:"limit,default=20"`
}

// 3. Response untuk menampilkan data Word ke Frontend
type WordResponse struct {
	ID                 string           `json:"id"`
	TargetLanguageCode string           `json:"target_language_code"`
	TargetWord         string           `json:"target_word"`
	NativeWord         string           `json:"native_word"`
	PartOfSpeech       string           `json:"part_of_speech"`
	IsFavorite         bool             `json:"is_favorite"`
	IsBookmarked       bool             `json:"is_bookmarked"`
	CreatedAt          string           `json:"created_at"`
	Categories         []CategoryRes    `json:"categories"`
	Examples           []WordExampleRes `json:"examples"`
}

type CategoryRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type WordExampleRes struct {
	ID             string `json:"id"`
	TargetSentence string `json:"target_sentence"`
	NativeSentence string `json:"native_sentence"`
}

// Tambahkan struct ini di bawah WordFilterRequest

type PaginationMeta struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
}

type WordPaginationResponse struct {
	Data []WordResponse `json:"data"` // Data tetap menggunakan struktur list kata yang lama
	Meta PaginationMeta `json:"meta"`
}

// Tambahkan struct ini untuk respons Part of Speech
type WordTypeCountResponse struct {
	PartOfSpeech string `json:"part_of_speech"`
	Count        int64  `json:"count"`
}

// DTO Khusus untuk Update (Hanya 1 bahasa karena merujuk pada 1 ID Word spesifik)
type UpdateWordRequest struct {
	TargetWord     string           `json:"target_word" binding:"required"`
	NativeWord     string           `json:"native_word" binding:"required"`
	PartOfSpeech   string           `json:"part_of_speech" binding:"required"`
	CategoryIDs    []string         `json:"category_ids"`
	Examples       []ExampleRequest `json:"examples" binding:"max=3"`
}