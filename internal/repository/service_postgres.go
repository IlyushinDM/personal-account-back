package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// ServicePostgres реализует ServiceRepository для PostgreSQL.
type ServicePostgres struct {
	db *sqlx.DB
}

// NewServicePostgres создает новый экземпляр репозитория для услуг.
func NewServicePostgres(db *sqlx.DB) *ServicePostgres {
	return &ServicePostgres{db: db}
}

// GetServiceRecommendations возвращает фиктивный текст рекомендаций для услуги.
// ! Это mock-реализация. В реальном приложении здесь будет запрос к полю в таблице services.
func (s *ServicePostgres) GetServiceRecommendations(ctx context.Context, serviceID uint64) (string, error) {
	// В реальном приложении здесь была бы проверка существования serviceID
	if serviceID == 0 {
		return "", sql.ErrNoRows // Используем стандартную ошибку для "не найдено"
	}
	return "Рекомендуется не принимать пищу за 4 часа до процедуры. Воду пить можно.", nil
}
