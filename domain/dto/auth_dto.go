package dto

// Format data yang dikirim user saat Register
type RegisterRequest struct {
	Name             string               `json:"name" binding:"required"`
	Username         string               `json:"username" binding:"required"`
	Email            string               `json:"email" binding:"required,email"`
	Password         string               `json:"password" binding:"required,min=6"`
	NoWa             string               `json:"no_wa"`
	NativeLanguage   string               `json:"native_language" binding:"required"`
	NativeFlagIcon   string               `json:"native_flag_icon" binding:"required"`
	Slogan           string               `json:"slogan"`
	TargetLanguages  []TargetLanguageReq  `json:"target_languages" binding:"required,min=1,max=2"` // Maksimal 2 bahasa
}

type TargetLanguageReq struct {
	LanguageCode string `json:"language_code" binding:"required"`
	LanguageName string `json:"language_name" binding:"required"`
	FlagIcon     string `json:"flag_icon" binding:"required"`
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
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	Username           string               `json:"username"`
	Email              string               `json:"email"`
	SubscriptionStatus string               `json:"subscription_status"`
	NativeLanguage     string               `json:"native_language"`
	NativeFlagIcon     string               `json:"native_flag_icon"`
	Slogan             string               `json:"slogan"`
	TargetLanguages    []TargetLanguageRes  `json:"target_languages"`
}

type TargetLanguageRes struct {
	LanguageCode string `json:"language_code"`
	LanguageName string `json:"language_name"`
	FlagIcon     string `json:"flag_icon"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Tambahkan struct ini di bagian bawah
type UpdateProfileRequest struct {
	Name             string               `json:"name" binding:"required"`
	Username         string               `json:"username" binding:"required"`
	NoWa             string               `json:"no_wa"`
	NativeLanguage   string               `json:"native_language" binding:"required"`
	NativeFlagIcon   string               `json:"native_flag_icon" binding:"required"`
	Slogan           string               `json:"slogan"`
	TargetLanguages  []TargetLanguageReq  `json:"target_languages" binding:"required,min=1,max=2"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	ResetToken  string `json:"reset_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}