package repository

import (
	"context"
	"pamagi/domain"
	"pamagi/domain/entity"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Create(c context.Context, user *entity.User) error {
	return r.db.WithContext(c).Create(user).Error
}

// Tambahkan di interface domain/auth.go dulu: GetByID(c context.Context, id string) (entity.User, error)
func (r *authRepository) GetByID(c context.Context, id string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(c).Preload("TargetLanguages").Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *authRepository) GetByEmail(c context.Context, email string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(c).Preload("TargetLanguages").Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *authRepository) GetByUsername(c context.Context, username string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(c).Preload("TargetLanguages").Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *authRepository) GetByIdentifier(c context.Context, identifier string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(c).Preload("TargetLanguages").Where("email = ? OR username = ?", identifier, identifier).First(&user).Error
	return user, err
}