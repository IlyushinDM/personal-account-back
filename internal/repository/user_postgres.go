package repository

import (
	"context"
	"fmt"

	"lk/internal/models"

	"gorm.io/gorm"
)

// UserPostgres реализует UserRepository для PostgreSQL.
type UserPostgres struct {
	db *gorm.DB
}

// NewUserPostgres создает новый экземпляр репозитория для пользователей.
func NewUserPostgres(db *gorm.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

// CreateUser вставляет нового пользователя в таблицу users и возвращает его ID.
// * Эта функция должна вызываться внутри транзакции.
func (r *UserPostgres) CreateUser(ctx context.Context, tx *gorm.DB, user models.User) (uint64, error) {
	result := tx.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

// CreateUserProfile вставляет новый профиль пользователя и возвращает его ID.
// * Эта функция должна вызываться внутри транзакции.
func (r *UserPostgres) CreateUserProfile(ctx context.Context, tx *gorm.DB, profile models.UserProfile) (uint64, error) {
	result := tx.WithContext(ctx).Create(&profile)
	if result.Error != nil {
		return 0, result.Error
	}
	return profile.ID, nil
}

// GetUserByPhone находит пользователя по номеру телефона.
func (r *UserPostgres) GetUserByPhone(ctx context.Context, phone string) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

// GetUserByID находит пользователя по его ID.
func (r *UserPostgres) GetUserByID(ctx context.Context, id uint64) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return user, err
}

// GetUserProfileByUserID находит профиль пользователя по ID пользователя.
func (r *UserPostgres) GetUserProfileByUserID(ctx context.Context, userID uint64) (models.UserProfile, error) {
	var profile models.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	return profile, err
}

// UpdateUserProfile обновляет данные профиля пользователя.
func (r *UserPostgres) UpdateUserProfile(ctx context.Context, profile models.UserProfile) (models.UserProfile, error) {
	updateData := make(map[string]interface{})

	if profile.Email.Valid {
		updateData["email"] = profile.Email
	}
	if profile.CityID > 0 {
		updateData["city_id"] = profile.CityID
	}

	if len(updateData) == 0 {
		// Если нечего обновлять, просто возвращаем текущий профиль
		return r.GetUserProfileByUserID(ctx, profile.UserID)
	}

	result := r.db.WithContext(ctx).Model(&models.UserProfile{}).Where(
		"user_id = ?", profile.UserID).Updates(updateData)
	if result.Error != nil {
		return models.UserProfile{}, result.Error
	}

	// Возвращаем обновленный профиль
	var updatedProfile models.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", profile.UserID).First(&updatedProfile).Error
	return updatedProfile, err
}

// UpdateAvatar обновляет URL аватара для профиля пользователя.
func (r *UserPostgres) UpdateAvatar(ctx context.Context, userID uint64, avatarURL string) error {
	result := r.db.WithContext(ctx).Model(&models.UserProfile{}).Where(
		"user_id = ?", userID).Update("avatar_url", avatarURL)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user profile with user_id %d not found for avatar update", userID)
	}
	return nil
}

// UpdatePassword обновляет хеш пароля пользователя.
func (r *UserPostgres) UpdatePassword(ctx context.Context, userID uint64, newPasswordHash string) error {
	result := r.db.WithContext(ctx).Model(&models.User{}).Where(
		"id = ?", userID).Update("password_hash", newPasswordHash)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id %d not found for password update", userID)
	}
	return nil
}
