package repository

import (
	"context"

	"lk/internal/models"

	"gorm.io/gorm"
)

// ServicePostgres реализует ServiceRepository для PostgreSQL.
type ServicePostgres struct {
	db *gorm.DB
}

// NewServicePostgres создает новый экземпляр репозитория.
func NewServicePostgres(db *gorm.DB) *ServicePostgres {
	return &ServicePostgres{db: db}
}

// GetServiceRecommendations получает рекомендации из поля в таблице services.
func (s *ServicePostgres) GetServiceRecommendations(ctx context.Context, serviceID uint64) (string, error) {
	var service models.Service
	err := s.db.WithContext(ctx).Select("recommendations").First(&service, serviceID).Error
	if err != nil {
		return "", err
	}
	return service.Recommendations.String, nil
}
