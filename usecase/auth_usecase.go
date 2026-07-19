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
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	authRepo domain.AuthRepository
	env      *bootstrap.Env
	redis    *redis.Client
}

func NewAuthUsecase(authRepo domain.AuthRepository, env *bootstrap.Env, redis *redis.Client) domain.AuthUsecase {
	return &authUsecase{
		authRepo: authRepo,
		env:      env,
		redis:    redis,
	}
}

func (u *authUsecase) Register(c context.Context, req *dto.RegisterRequest) (dto.AuthResponse, error) {
	if _, err := u.authRepo.GetByEmail(c, req.Email); err == nil {
		return dto.AuthResponse{}, errors.New("email sudah terdaftar")
	}
	if _, err := u.authRepo.GetByUsername(c, req.Username); err == nil {
		return dto.AuthResponse{}, errors.New("username sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	userId := uuid.New().String()
	
	// Handle Slogan Default
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

	// Mapping Target Languages
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

	// GORM akan otomatis menyimpan User dan TargetLanguages (relasi Has-Many) secara bersamaan
	if err := u.authRepo.Create(c, user); err != nil {
		return dto.AuthResponse{}, err
	}

	// Berikan seluruh entity user ke helper
	return u.generateTokensAndStore(c, user)
}

func (u *authUsecase) Login(c context.Context, req *dto.LoginRequest) (dto.AuthResponse, error) {
	// Ubah GetByEmail menjadi GetByIdentifier dan panggil req.Identifier
	user, err := u.authRepo.GetByIdentifier(c, req.Identifier)
	if err != nil {
		return dto.AuthResponse{}, errors.New("kredensial tidak valid")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return dto.AuthResponse{}, errors.New("kredensial tidak valid")
	}

	return u.generateTokensAndStore(c, &user)
}

// Helper internal untuk generate token dan simpan ke Redis
func (u *authUsecase) generateTokensAndStore(c context.Context, user *entity.User) (dto.AuthResponse, error) {
	accessToken, err := tokenutil.CreateAccessToken(user.ID, user.Name, u.env.AccessTokenSecret, 15)
	if err != nil { return dto.AuthResponse{}, err }

	refreshToken, err := tokenutil.CreateRefreshToken(user.ID, u.env.AccessTokenSecret, 7)
	if err != nil { return dto.AuthResponse{}, err }

	redisKey := "refresh_token:" + user.ID
	duration := time.Hour * 24 * 7 
	if err := u.redis.Set(c, redisKey, refreshToken, duration).Err(); err != nil {
		return dto.AuthResponse{}, errors.New("gagal menyimpan sesi di server")
	}

	// Mapping relasi bahasa ke DTO Response
	var targetRes []dto.TargetLanguageRes
	for _, t := range user.TargetLanguages {
		targetRes = append(targetRes, dto.TargetLanguageRes{
			LanguageCode: t.LanguageCode,
			LanguageName: t.LanguageName,
			FlagIcon:     t.FlagIcon,
		})
	}

	// Handle pointer untuk safety
	nativeLang := ""
	nativeIcon := ""
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
	// Kunci token di Redis sesuai dengan format saat login
	redisKey := "refresh_token:" + userID
	
	// Hapus token dari Redis
	return u.redis.Del(c, redisKey).Err()
}

// Tambahkan di interface domain/auth.go: GetProfile(c context.Context, userID string) (dto.UserResponse, error)
func (u *authUsecase) GetProfile(c context.Context, userID string) (dto.UserResponse, error) {
	user, err := u.authRepo.GetByID(c, userID)
	if err != nil {
		return dto.UserResponse{}, err
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
	// 1. Validasi keaslian token
	authorized, err := tokenutil.IsAuthorized(req.RefreshToken, u.env.AccessTokenSecret)
	if !authorized || err != nil {
		return dto.AuthResponse{}, errors.New("refresh token tidak valid")
	}

	// 2. Ekstrak userID dari dalam token
	userID, err := tokenutil.ExtractIDFromToken(req.RefreshToken, u.env.AccessTokenSecret)
	if err != nil {
		return dto.AuthResponse{}, errors.New("gagal membaca token")
	}

	// 3. Cek di Redis apakah token tersebut cocok dan belum di-logout
	redisKey := "refresh_token:" + userID
	storedToken, err := u.redis.Get(c, redisKey).Result()
	if err != nil || storedToken != req.RefreshToken {
		return dto.AuthResponse{}, errors.New("sesi telah kedaluwarsa atau tidak valid, silakan login kembali")
	}

	// 4. Ambil data user lengkap dari database
	user, err := u.authRepo.GetByID(c, userID)
	if err != nil {
		return dto.AuthResponse{}, errors.New("user tidak ditemukan")
	}

	// 5. Terbitkan pasangan token baru (Access & Refresh) lalu simpan ke Redis
	// Kita manfaatkan helper yang sudah ada
	return u.generateTokensAndStore(c, &user)
}