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
	user := &entity.User{
		ID:           userId,
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		NoWa:         req.NoWa,
	}

	if req.Region != "" {
		user.Region = &req.Region
	}

	if err := u.authRepo.Create(c, user); err != nil {
		return dto.AuthResponse{}, err
	}

	// Setelah sukses buat user, langsung panggil logika login otomatis
	return u.generateTokensAndStore(c, userId, user.Name, user.Username, user.Email, user.SubscriptionStatus)
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

	return u.generateTokensAndStore(c, user.ID, user.Name, user.Username, user.Email, user.SubscriptionStatus)
}

// Helper internal untuk generate token dan simpan ke Redis
func (u *authUsecase) generateTokensAndStore(c context.Context, userId, name, username, email string) (dto.AuthResponse, error) {
	// 1. Buat Access Token (15 Menit)
	accessToken, err := tokenutil.CreateAccessToken(userId, name, u.env.AccessTokenSecret, 15)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	// 2. Buat Refresh Token (7 Hari)
	refreshToken, err := tokenutil.CreateRefreshToken(userId, u.env.AccessTokenSecret, 7)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	// 3. Simpan Refresh Token ke Redis (Key -> refresh_token:user_id)
	redisKey := "refresh_token:" + userId
	duration := time.Hour * 24 * 7 // 7 Hari
	if err := u.redis.Set(c, redisKey, refreshToken, duration).Err(); err != nil {
		return dto.AuthResponse{}, errors.New("gagal menyimpan sesi di server")
	}

	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       userId,
			Name:     name,
			Username: username,
			Email:    email,
			SubscriptionStatus: subscriptionStatus,
		},
	}, nil
}


func (u *authUsecase) Logout(c context.Context, userID string) error {
	// Kunci token di Redis sesuai dengan format saat login
	redisKey := "refresh_token:" + userID
	
	// Hapus token dari Redis
	return u.redis.Del(c, redisKey).Err()
}