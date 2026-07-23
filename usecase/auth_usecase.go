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
	accessToken, err := tokenutil.CreateAccessToken(user.ID, user.Name, u.env.AccessTokenSecret, 15)
	if err != nil { return dto.AuthResponse{}, err }

	refreshToken, err := tokenutil.CreateRefreshToken(user.ID, u.env.AccessTokenSecret, 7)
	if err != nil { return dto.AuthResponse{}, err }

	redisKey := "refresh_token:" + user.ID
	duration := time.Hour * 24 * 7 
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