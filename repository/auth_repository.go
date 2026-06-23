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

func (r *authRepository) GetByEmail(c context.Context, email string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(c).Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *authRepository) GetByUsername(c context.Context, username string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(c).Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *authRepository) GetByIdentifier(c context.Context, identifier string) (entity.User, error) {
	var user entity.User
	// GORM akan mencari: WHERE email = 'identifier' OR username = 'identifier'
	err := r.db.WithContext(c).Where("email = ? OR username = ?", identifier, identifier).First(&user).Error
	return user, err
}