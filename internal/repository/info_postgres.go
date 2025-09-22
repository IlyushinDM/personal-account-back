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

// NewInfoPostgres создает новый экземпляр репозитория для общей информации.
func NewInfoPostgres(db *gorm.DB) *InfoPostgres {
	return &InfoPostgres{db: db}
}

// GetClinicInfo возвращает mock-данные о клинике.
// ! Это mock-реализация. В реальном приложении здесь будет запрос к БД.
func (r *InfoPostgres) GetClinicInfo(ctx context.Context) (models.ClinicInfo, error) {
	info := models.ClinicInfo{
		Name: "Клиника 'Здоровье'",
		Contacts: []models.Contact{
			{Type: "phone", Value: "+7 (495) 222-33-44"},
			{Type: "email", Value: "info@clinic.ru"},
		},
		Addresses: []models.Address{
			{ID: 1, Address: "ул. Ленина, д. 100", IsMain: true},
		},
		WorkingHours: []models.WorkHours{
			{Days: "пн-пт", Hours: "08:00 - 20:00"},
			{Days: "сб", Hours: "09:00 - 18:00"},
		},
	}
	// В реальном приложении можно было бы сделать r.db.WithContext(ctx).Find(&info)
	return info, nil
}

// GetLegalDocuments возвращает mock-данные о юридических документах.
// ! Это mock-реализация.
func (r *InfoPostgres) GetLegalDocuments(ctx context.Context) ([]models.LegalDocument, error) {
	docs := []models.LegalDocument{
		{
			Type:       "terms_of_use",
			Title:      "Пользовательское соглашение",
			URL:        "/legal/terms.pdf",
			Version:    "1.2",
			UpdateDate: "2023-10-01",
		},
		{
			Type:       "privacy_policy",
			Title:      "Политика конфиденциальности",
			URL:        "/legal/privacy.pdf",
			Version:    "2.0",
			UpdateDate: "2023-09-15",
		},
	}
	return docs, nil
}
