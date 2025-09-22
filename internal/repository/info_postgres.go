package repository

import (
	"context"

	"lk/internal/models"

	"gorm.io/gorm"
)

// InfoPostgres реализует InfoRepository.
type InfoPostgres struct {
	db *gorm.DB
}

// NewInfoPostgres создает новый экземпляр репозитория.
func NewInfoPostgres(db *gorm.DB) *InfoPostgres {
	return &InfoPostgres{db: db}
}

// GetClinicInfo возвращает информацию о клинике из БД.
func (r *InfoPostgres) GetClinicInfo(ctx context.Context) (models.ClinicInfo, error) {
	var clinic models.Clinic
	// Предполагаем, что основная клиника имеет ID = 1
	err := r.db.WithContext(ctx).First(&clinic, 1).Error
	if err != nil {
		return models.ClinicInfo{}, err
	}

	// Формируем DTO на основе полученных данных
	info := models.ClinicInfo{
		Name: clinic.Name,
		Contacts: []models.Contact{
			{Type: "phone", Value: clinic.Phone},
			{Type: "email", Value: "info@clinic.ru"},
		},
		Addresses: []models.Address{
			{ID: int(clinic.ID), Address: clinic.Address, IsMain: true},
		},
		WorkingHours: []models.WorkHours{
			{Days: "пн-пт", Hours: "08:00 - 20:00"},
			{Days: "сб", Hours: "09:00 - 18:00"},
		},
	}
	return info, nil
}

// GetLegalDocuments возвращает список юридических документов из БД.
func (r *InfoPostgres) GetLegalDocuments(ctx context.Context) ([]models.LegalDocument, error) {
	var docs []models.LegalDocument
	err := r.db.WithContext(ctx).Find(&docs).Error
	return docs, err
}
