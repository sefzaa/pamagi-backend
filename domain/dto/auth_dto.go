package dto

// Format data yang dikirim user saat Register
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	NoWa     string `json:"no_wa"`
	Region   string `json:"region" binding:"omitempty,iso3166_1_alpha2"`
}

// Format data yang dikirim user saat Login
type LoginRequest struct {
	// Ubah Email menjadi Identifier (bisa menampung email atau username)
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// Format balasan sukses setelah Login
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// Format data user yang aman dikembalikan (tanpa password)
type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	SubscriptionStatus string `json:"subscription_status"`
}