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
}

type AuthUsecase interface {
	Register(c context.Context, req *dto.RegisterRequest) (dto.AuthResponse, error)
	Login(c context.Context, req *dto.LoginRequest) (dto.AuthResponse, error)
	Logout(c context.Context, userID string) error
	GetProfile(c context.Context, userID string) (dto.UserResponse, error) // TAMBAHAN INI
	RefreshToken(c context.Context, req *dto.RefreshTokenRequest) (dto.AuthResponse, error)
}