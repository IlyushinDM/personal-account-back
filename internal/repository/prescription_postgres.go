package repository

import (
	"context"
	"errors"
	"time"

	"lk/internal/models"

	"gorm.io/gorm"
)

// PrescriptionPostgres реализует PrescriptionRepository.
type PrescriptionPostgres struct {
	db *gorm.DB
}

// NewPrescriptionPostgres создает новый экземпляр репозитория для назначений.
func NewPrescriptionPostgres(db *gorm.DB) *PrescriptionPostgres {
	return &PrescriptionPostgres{db: db}
}

// GetActiveByUserID возвращает mock-данные об активных назначениях.
// ! Это mock-реализация.
func (r *PrescriptionPostgres) GetActiveByUserID(ctx context.Context, userID uint64) ([]models.Prescription, error) {
	// В реальном приложении запрос выглядел бы так:
	// var prescriptions []models.Prescription
	// err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, "active").Find(&prescriptions).Error
	// return prescriptions, err

	prescriptions := []models.Prescription{
		{
			ID:            888,
			AppointmentID: 554,
			UserID:        userID,
			DoctorID:      77,
			Content:       "УЗИ брюшной полости",
			Status:        "active",
			CreatedAt:     time.Now().Add(-5 * 24 * time.Hour),
		},
	}
	return prescriptions, nil
}

// Archive помечает назначение как архивированное.
// ! Это mock-реализация.
func (r *PrescriptionPostgres) Archive(ctx context.Context, userID, prescriptionID uint64) error {
	// В реальном приложении здесь будет обновление статуса в БД:
	// err := r.db.WithContext(ctx).Model(&models.Prescription{}).
	//    Where("id = ? AND user_id = ?", prescriptionID, userID).
	//    Update("status", "archived").Error
	// return err

	if prescriptionID == 0 {
		return errors.New("prescription not found")
	}
	return nil
}

// GetByID возвращает назначение по ID.
// ! Это mock-реализация.
func (r *PrescriptionPostgres) GetByID(ctx context.Context, prescriptionID uint64) (models.Prescription, error) {
	// В реальном приложении:
	// var prescription models.Prescription
	// err := r.db.WithContext(ctx).First(&prescription, prescriptionID).Error
	// return prescription, err

	if prescriptionID != 888 {
		return models.Prescription{}, gorm.ErrRecordNotFound
	}
	return models.Prescription{
		ID:            888,
		AppointmentID: 554,
		UserID:        1, // Mock UserID
		DoctorID:      77,
		Content:       "УЗИ брюшной полости",
		Status:        "active",
		CreatedAt:     time.Now().Add(-5 * 24 * time.Hour),
	}, nil
}
