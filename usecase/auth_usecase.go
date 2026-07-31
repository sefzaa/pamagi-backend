package usecase

import (
	"context"
	"errors"
	"time"
	"pamagi/bootstrap"
	"pamagi/domain"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
	"pamagi/internal/tokenutil"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/redis/go-redis/v9"

	"crypto/rand"
	"fmt"
	"math/big"
)

type authUsecase struct {
	authRepo domain.AuthRepository
	env      *bootstrap.Env
	redis    *redis.Client
	emailService domain.EmailService
}

func NewAuthUsecase(authRepo domain.AuthRepository, env *bootstrap.Env, redis *redis.Client, emailService domain.EmailService) domain.AuthUsecase {
	return &authUsecase{
		authRepo: authRepo,
		env:      env,
		redis:    redis,
		emailService: emailService,
	}
}

func (u *authUsecase) Register(c context.Context, req *dto.RegisterRequest) (dto.AuthResponse, error) {
	if _, err := u.authRepo.GetByEmail(c, req.Email); err == nil {
		return dto.AuthResponse{}, errors.New("Email is already registered")
	}
	if _, err := u.authRepo.GetByUsername(c, req.Username); err == nil {
		return dto.AuthResponse{}, errors.New("Username is already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	userId := uuid.New().String()
	
	slogan := req.Slogan
	if slogan == "" {
		slogan = "Consistency is key to fluency."
	}

	user := &entity.User{
		ID:             userId,
		Name:           req.Name,
		Username:       req.Username,
		Email:          req.Email,
		PasswordHash:   string(hashedPassword),
		NoWa:           req.NoWa,
		NativeLanguage: &req.NativeLanguage,
		NativeFlagIcon: &req.NativeFlagIcon,
		Slogan:         slogan,
	}

	var targetLanguages []entity.UserTargetLanguage
	for _, targetReq := range req.TargetLanguages {
		targetLanguages = append(targetLanguages, entity.UserTargetLanguage{
			ID:           uuid.New().String(),
			UserID:       userId,
			LanguageCode: targetReq.LanguageCode,
			LanguageName: targetReq.LanguageName,
			FlagIcon:     targetReq.FlagIcon,
		})
	}
	user.TargetLanguages = targetLanguages

	if err := u.authRepo.Create(c, user); err != nil {
		return dto.AuthResponse{}, err
	}

	return u.generateTokensAndStore(c, user)
}

func (u *authUsecase) Login(c context.Context, req *dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := u.authRepo.GetByIdentifier(c, req.Identifier)
	if err != nil {
		return dto.AuthResponse{}, errors.New("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return dto.AuthResponse{}, errors.New("Invalid credentials")
	}

	return u.generateTokensAndStore(c, &user)
}

func (u *authUsecase) generateTokensAndStore(c context.Context, user *entity.User) (dto.AuthResponse, error) {
	accessToken, err := tokenutil.CreateAccessToken(user.ID, user.Name, u.env.AccessTokenSecret, 15) // Access token tetap pendek (15 menit)
	if err != nil { return dto.AuthResponse{}, err }

	// UBAH: Set umur Refresh Token menjadi 30 Hari
	refreshToken, err := tokenutil.CreateRefreshToken(user.ID, u.env.AccessTokenSecret, 30)
	if err != nil { return dto.AuthResponse{}, err }

	redisKey := "refresh_token:" + user.ID
	// UBAH: Set durasi penyimpanan di Redis menjadi 30 Hari
	duration := time.Hour * 24 * 30 
	
	if err := u.redis.Set(c, redisKey, refreshToken, duration).Err(); err != nil {
		return dto.AuthResponse{}, errors.New("Failed to save session on server")
	}

	var targetRes []dto.TargetLanguageRes
	for _, t := range user.TargetLanguages {
		targetRes = append(targetRes, dto.TargetLanguageRes{
			LanguageCode: t.LanguageCode,
			LanguageName: t.LanguageName,
			FlagIcon:     t.FlagIcon,
		})
	}

	nativeLang, nativeIcon := "", ""
	if user.NativeLanguage != nil { nativeLang = *user.NativeLanguage }
	if user.NativeFlagIcon != nil { nativeIcon = *user.NativeFlagIcon }

	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:                 user.ID,
			Name:               user.Name,
			Username:           user.Username,
			Email:              user.Email,
			SubscriptionStatus: user.SubscriptionStatus,
			NativeLanguage:     nativeLang,
			NativeFlagIcon:     nativeIcon,
			Slogan:             user.Slogan,
			TargetLanguages:    targetRes,
		},
	}, nil
}

func (u *authUsecase) Logout(c context.Context, userID string) error {
	redisKey := "refresh_token:" + userID
	return u.redis.Del(c, redisKey).Err()
}

func (u *authUsecase) GetProfile(c context.Context, userID string) (dto.UserResponse, error) {
	user, err := u.authRepo.GetByID(c, userID)
	if err != nil {
		return dto.UserResponse{}, errors.New("User not found")
	}

	var targetRes []dto.TargetLanguageRes
	for _, t := range user.TargetLanguages {
		targetRes = append(targetRes, dto.TargetLanguageRes{
			LanguageCode: t.LanguageCode,
			LanguageName: t.LanguageName,
			FlagIcon:     t.FlagIcon,
		})
	}

	nativeLang, nativeIcon := "", ""
	if user.NativeLanguage != nil { nativeLang = *user.NativeLanguage }
	if user.NativeFlagIcon != nil { nativeIcon = *user.NativeFlagIcon }

	return dto.UserResponse{
		ID:                 user.ID,
		Name:               user.Name,
		Username:           user.Username,
		Email:              user.Email,
		SubscriptionStatus: user.SubscriptionStatus,
		NativeLanguage:     nativeLang,
		NativeFlagIcon:     nativeIcon,
		Slogan:             user.Slogan,
		TargetLanguages:    targetRes,
	}, nil
}

func (u *authUsecase) RefreshToken(c context.Context, req *dto.RefreshTokenRequest) (dto.AuthResponse, error) {
	authorized, err := tokenutil.IsAuthorized(req.RefreshToken, u.env.AccessTokenSecret)
	if !authorized || err != nil {
		return dto.AuthResponse{}, errors.New("Invalid refresh token")
	}

	userID, err := tokenutil.ExtractIDFromToken(req.RefreshToken, u.env.AccessTokenSecret)
	if err != nil {
		return dto.AuthResponse{}, errors.New("Failed to parse token")
	}

	redisKey := "refresh_token:" + userID
	storedToken, err := u.redis.Get(c, redisKey).Result()
	if err != nil || storedToken != req.RefreshToken {
		return dto.AuthResponse{}, errors.New("Session expired, please login again")
	}

	user, err := u.authRepo.GetByID(c, userID)
	if err != nil {
		return dto.AuthResponse{}, errors.New("User not found")
	}

	return u.generateTokensAndStore(c, &user)
}

func (u *authUsecase) UpdateProfile(c context.Context, userID string, req *dto.UpdateProfileRequest) (dto.UserResponse, error) {
	// Cek apakah username sudah dipakai orang lain
	existingUser, err := u.authRepo.GetByUsername(c, req.Username)
	if err == nil && existingUser.ID != userID {
		return dto.UserResponse{}, errors.New("Username is already taken by someone else")
	}

	user := &entity.User{
		ID:             userID,
		Name:           req.Name,
		Username:       req.Username,
		NoWa:           req.NoWa,
		NativeLanguage: &req.NativeLanguage,
		NativeFlagIcon: &req.NativeFlagIcon,
		Slogan:         req.Slogan,
	}

	var targetLanguages []entity.UserTargetLanguage
	for _, targetReq := range req.TargetLanguages {
		targetLanguages = append(targetLanguages, entity.UserTargetLanguage{
			ID:           uuid.New().String(),
			UserID:       userID,
			LanguageCode: targetReq.LanguageCode,
			LanguageName: targetReq.LanguageName,
			FlagIcon:     targetReq.FlagIcon,
		})
	}
	user.TargetLanguages = targetLanguages

	// Eksekusi update
	if err := u.authRepo.UpdateProfile(c, user); err != nil {
		return dto.UserResponse{}, errors.New("Failed to update profile")
	}

	// Kembalikan profil yang baru dengan menggunakan fungsi GetProfile yang sudah ada
	return u.GetProfile(c, userID)
}

func generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// Implementasi Forgot Password
func (u *authUsecase) ForgotPassword(c context.Context, req *dto.ForgotPasswordRequest) error {
	// 1. Cek apakah email ada di database
	user, err := u.authRepo.GetByEmail(c, req.Email)
	if err != nil {
		// Untuk keamanan, tetap kembalikan nil agar peretas tidak tahu email terdaftar atau tidak
		return nil 
	}

	// 2. Buat OTP
	otpCode, err := generateOTP()
	if err != nil {
		return errors.New("gagal membuat OTP")
	}

	// 3. Simpan OTP ke Redis dengan batas waktu 5 menit
	redisKey := "otp_forgot_password:" + user.Email
	err = u.redis.Set(c, redisKey, otpCode, 5*time.Minute).Err()
	if err != nil {
		return errors.New("gagal menyimpan OTP ke sistem")
	}

	// 4. Kirim email via Brevo
	go func() {
		// Gunakan goroutine agar respons API tidak menunggu email terkirim (asynchronous)
		bgCtx := context.Background() 
		_ = u.emailService.SendOTP(bgCtx, user.Email, otpCode)
	}()

	return nil
}

func (u *authUsecase) VerifyOTP(c context.Context, req *dto.VerifyOTPRequest) (string, error) {
	redisKey := "otp_forgot_password:" + req.Email
	storedOTP, err := u.redis.Get(c, redisKey).Result()
	
	if err != nil || storedOTP != req.OTP {
		return "", errors.New("Kode OTP tidak valid atau sudah kedaluwarsa")
	}

	// 1. Jika OTP benar, HAPUS OTP dari Redis agar tidak bisa dipakai ulang
	u.redis.Del(c, redisKey)

	// 2. Buat "Reset Token" sebagai tiket masuk ke endpoint ResetPassword
	resetToken := uuid.New().String()
	resetKey := "reset_token:" + req.Email
	
	// Simpan tiket ini di Redis dengan masa berlaku 15 menit
	err = u.redis.Set(c, resetKey, resetToken, 15*time.Minute).Err()
	if err != nil {
		return "", errors.New("Gagal membuat sesi pemulihan")
	}

	return resetToken, nil
}

func (u *authUsecase) ResetPassword(c context.Context, req *dto.ResetPasswordRequest) error {
	resetKey := "reset_token:" + req.Email
	storedToken, err := u.redis.Get(c, resetKey).Result()

	// Cek apakah tiket (Reset Token) valid
	if err != nil || storedToken != req.ResetToken {
		return errors.New("Sesi reset password tidak valid atau sudah habis. Silakan ulangi dari awal.")
	}

	// Ambil data user
	user, err := u.authRepo.GetByEmail(c, req.Email)
	if err != nil {
		return errors.New("User tidak ditemukan")
	}

	// Hash password yang baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Gagal memproses password baru")
	}

	// Simpan ke database
	if err := u.authRepo.UpdatePassword(c, user.ID, string(hashedPassword)); err != nil {
		return errors.New("Gagal mengubah password di database")
	}

	// Hapus Reset Token agar tiket tidak bisa dipakai lagi (One-time use)
	u.redis.Del(c, resetKey)
	
	// Opsional: Hapus juga sesi refresh token (logout paksa semua perangkat)
	_ = u.redis.Del(c, "refresh_token:"+user.ID)

	return nil
}