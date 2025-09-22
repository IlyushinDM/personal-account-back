package repository

import (
	"context"

	"gorm.io/gorm"
)

// ServicePostgres реализует ServiceRepository для PostgreSQL.
type ServicePostgres struct {
	db *gorm.DB
}

// NewServicePostgres создает новый экземпляр репозитория для услуг.
func NewServicePostgres(db *gorm.DB) *ServicePostgres {
	return &ServicePostgres{db: db}
}

// GetServiceRecommendations возвращает фиктивный текст рекомендаций для услуги.
// ! Это mock-реализация. В реальном приложении здесь будет запрос к полю в таблице services.
func (s *ServicePostgres) GetServiceRecommendations(ctx context.Context, serviceID uint64) (string, error) {
	// В реальном приложении здесь была бы проверка существования serviceID
	if serviceID == 0 {
		return "", gorm.ErrRecordNotFound // Используем стандартную ошибку GORM для "не найдено"
	}
	return "Рекомендуется не принимать пищу за 4 часа до процедуры. Воду пить можно.", nil
}
