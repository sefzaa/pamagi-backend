package dto

// 1. Request untuk Create Word
type CreateWordRequest struct {
	RussianWord  string           `json:"russian_word" binding:"required"`
	Translation  string           `json:"translation" binding:"required"`
	PartOfSpeech string           `json:"part_of_speech" binding:"required"`
	CategoryIDs  []string         `json:"category_ids"` // Array ID kategori (Opsional)
	Examples     []ExampleRequest `json:"examples"`     // Array contoh kalimat (Opsional)
}

// Sub-request untuk contoh kalimat
type ExampleRequest struct {
	RussianSentence    string `json:"russian_sentence" binding:"required"`
	TranslatedSentence string `json:"translated_sentence" binding:"required"`
}

// 2. Request untuk Filter GET /words (Pake Query Params)
type WordFilterRequest struct {
	CategoryID   string `form:"category_id"`
	PartOfSpeech string `form:"part_of_speech"`
	IsFavorite   *bool  `form:"is_favorite"`
	IsBookmarked *bool  `form:"is_bookmarked"`
	StartDate    string `form:"start_date"` // Format YYYY-MM-DD
	EndDate      string `form:"end_date"`   // Format YYYY-MM-DD
	SortBy       string `form:"sort_by"`    // newest, oldest, a_z, z_a
	Page         int    `form:"page,default=1"`
	Limit        int    `form:"limit,default=20"` // Limit 20 kata
}

// 3. Response untuk menampilkan data Word ke Frontend
type WordResponse struct {
	ID           string           `json:"id"`
	RussianWord  string           `json:"russian_word"`
	Translation  string           `json:"translation"`
	PartOfSpeech string           `json:"part_of_speech"`
	IsFavorite   bool             `json:"is_favorite"`
	IsBookmarked bool             `json:"is_bookmarked"`
	CreatedAt    string           `json:"created_at"`
	Categories   []CategoryRes    `json:"categories"`
	Examples     []WordExampleRes `json:"examples"`
}

type CategoryRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type WordExampleRes struct {
	ID                 string `json:"id"`
	RussianSentence    string `json:"russian_sentence"`
	TranslatedSentence string `json:"translated_sentence"`
}