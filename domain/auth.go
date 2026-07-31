package domain

import (
	"context"
	"pamagi/domain/dto"
	"pamagi/domain/entity"
)

type AuthRepository interface {
	Create(c context.Context, user *entity.User) error
	GetByEmail(c context.Context, email string) (entity.User, error)
	GetByUsername(c context.Context, username string) (entity.User, error)
	GetByIdentifier(c context.Context, identifier string) (entity.User, error)
	GetByID(c context.Context, id string) (entity.User, error) // TAMBAHAN INI
	UpdateProfile(c context.Context, user *entity.User) error
	UpdatePassword(c context.Context, userID string, passwordHash string) error // TAMBAHAN
}

type AuthUsecase interface {
	Register(c context.Context, req *dto.RegisterRequest) (dto.AuthResponse, error)
	Login(c context.Context, req *dto.LoginRequest) (dto.AuthResponse, error)
	Logout(c context.Context, userID string) error
	GetProfile(c context.Context, userID string) (dto.UserResponse, error)
	RefreshToken(c context.Context, req *dto.RefreshTokenRequest) (dto.AuthResponse, error)
	UpdateProfile(c context.Context, userID string, req *dto.UpdateProfileRequest) (dto.UserResponse, error)
	ForgotPassword(c context.Context, req *dto.ForgotPasswordRequest) error 
	VerifyOTP(c context.Context, req *dto.VerifyOTPRequest) (string, error) // Mengembalikan Reset Token
	ResetPassword(c context.Context, req *dto.ResetPasswordRequest) error
}