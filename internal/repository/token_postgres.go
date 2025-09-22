package repository

import (
	"context"
	"time"

	"lk/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenPostgres struct {
	db *gorm.DB
}

func NewTokenPostgres(db *gorm.DB) *TokenPostgres {
	return &TokenPostgres{db: db}
}

// Create сохраняет или обновляет токен для пользователя.
func (r *TokenPostgres) Create(ctx context.Context, token models.RefreshToken) error {
	// Upsert: если токен для userID уже есть - обновить, если нет - создать.
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token_hash", "expires_at"}),
	}).Create(&token).Error
}

// GetByUserID находит токен по UserID.
func (r *TokenPostgres) GetByUserID(ctx context.Context, userID uint64) (models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.WithContext(ctx).Where(
		"user_id = ? AND expires_at > ?", userID, time.Now()).First(&token).Error
	return token, err
}

// Delete удаляет токен пользователя (для logout).
func (r *TokenPostgres) Delete(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}
