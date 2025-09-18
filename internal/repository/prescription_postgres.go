package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"lk/internal/models"

	"github.com/jmoiron/sqlx"
)

// PrescriptionPostgres реализует PrescriptionRepository.
type PrescriptionPostgres struct {
	db *sqlx.DB
}

// NewPrescriptionPostgres создает новый экземпляр репозитория для назначений.
func NewPrescriptionPostgres(db *sqlx.DB) *PrescriptionPostgres {
	return &PrescriptionPostgres{db: db}
}

// GetActiveByUserID возвращает mock-данные об активных назначениях.
// ! Это mock-реализация.
func (r *PrescriptionPostgres) GetActiveByUserID(ctx context.Context, userID uint64) ([]models.Prescription, error) {
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
	// В реальном приложении здесь будет проверка, что назначение принадлежит userID
	// и обновление статуса в БД.
	if prescriptionID == 0 {
		return errors.New("prescription not found")
	}
	return nil
}

// GetByID возвращает назначение по ID.
// ! Это mock-реализация.
func (r *PrescriptionPostgres) GetByID(ctx context.Context, prescriptionID uint64) (models.Prescription, error) {
	if prescriptionID != 888 {
		return models.Prescription{}, sql.ErrNoRows
	}
	return models.Prescription{
		ID:            888,
		AppointmentID: 554,
		UserID:        1,
		DoctorID:      77,
		Content:       "УЗИ брюшной полости",
		Status:        "active",
		CreatedAt:     time.Now().Add(-5 * 24 * time.Hour),
	}, nil
}
