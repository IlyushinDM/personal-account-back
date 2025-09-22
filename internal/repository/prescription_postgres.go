package repository

import (
	"context"

	"lk/internal/models"

	"gorm.io/gorm"
)

// PrescriptionPostgres реализует PrescriptionRepository.
type PrescriptionPostgres struct {
	db *gorm.DB
}

// NewPrescriptionPostgres создает новый экземпляр репозитория.
func NewPrescriptionPostgres(db *gorm.DB) *PrescriptionPostgres {
	return &PrescriptionPostgres{db: db}
}

// GetActiveByUserID возвращает активные назначения пользователя.
func (r *PrescriptionPostgres) GetActiveByUserID(ctx context.Context, userID uint64) ([]models.Prescription, error) {
	var prescriptions []models.Prescription
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, "active").
		Order("created_at DESC").
		Find(&prescriptions).Error
	return prescriptions, err
}

// Archive помечает назначение как архивированное.
func (r *PrescriptionPostgres) Archive(ctx context.Context, userID, prescriptionID uint64) error {
	result := r.db.WithContext(ctx).Model(&models.Prescription{}).
		Where("id = ? AND user_id = ?", prescriptionID, userID).
		Update("status", "archived")

	if result.Error != nil {
		return result.Error
	}
	// Можно добавить проверку result.RowsAffected == 0 для возврата ошибки "не найдено"
	return nil
}

// GetByID возвращает назначение по ID.
func (r *PrescriptionPostgres) GetByID(ctx context.Context, prescriptionID uint64) (models.Prescription, error) {
	var prescription models.Prescription
	err := r.db.WithContext(ctx).First(&prescription, prescriptionID).Error
	return prescription, err
}
