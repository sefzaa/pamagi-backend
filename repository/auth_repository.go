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

func (r *authRepository) UpdateProfile(c context.Context, user *entity.User) error {
	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// 1. Update data user utama
		if err := tx.Model(user).Where("id = ?", user.ID).Updates(map[string]interface{}{
			"name":             user.Name,
			"username":         user.Username,
			"no_wa":            user.NoWa,
			"native_language":  user.NativeLanguage,
			"native_flag_icon": user.NativeFlagIcon,
			"slogan":           user.Slogan,
		}).Error; err != nil {
			return err
		}

		// 2. Hapus target bahasa yang lama
		if err := tx.Where("user_id = ?", user.ID).Delete(&entity.UserTargetLanguage{}).Error; err != nil {
			return err
		}

		// 3. Masukkan target bahasa yang baru
		if len(user.TargetLanguages) > 0 {
			if err := tx.Create(&user.TargetLanguages).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *authRepository) UpdatePassword(c context.Context, userID string, passwordHash string) error {
	// Mengupdate kolom password_hash berdasarkan userID
	return r.db.WithContext(c).Model(&entity.User{}).Where("id = ?", userID).Update("password_hash", passwordHash).Error
}