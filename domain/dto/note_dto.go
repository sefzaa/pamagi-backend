package dto

// Request saat User membuat atau mengedit Note
type NoteRequest struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"` // String mentah dari editor FE
	Color      string   `json:"color"`
	IsPinned   bool     `json:"is_pinned"`
	IsFavorite bool     `json:"is_favorite"`
	Tags       []string `json:"tags"` // Array of string untuk nama tag (misal: ["grammar", "vocabulary"])
}

// Response saat User mengambil daftar Note
type NoteResponse struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Preview    string   `json:"preview"`
	Color      string   `json:"color"`
	IsPinned   bool     `json:"is_pinned"`
	IsFavorite bool     `json:"is_favorite"`
	Tags       []string `json:"tags"`
	UpdatedAt  string   `json:"updated_at"`
}

// Response ringan untuk list (hanya preview konten)
type NoteListResponse struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Preview    string   `json:"preview"` // Potongan pendek dari isi catatan
	Color      string   `json:"color"`
	IsPinned   bool     `json:"is_pinned"`
	IsFavorite bool     `json:"is_favorite"`
	Tags       []string `json:"tags"`
	UpdatedAt  string   `json:"updated_at"` // Sesuai permintaan untuk tampil di card
}

// Response detail lengkap
type NoteDetailResponse struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"` // Isi utuh (Rich text / Delta)
	Color      string   `json:"color"`
	IsPinned   bool     `json:"is_pinned"`
	IsFavorite bool     `json:"is_favorite"`
	Tags       []string `json:"tags"`
	UpdatedAt  string   `json:"updated_at"`
}