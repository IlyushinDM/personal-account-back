package services

import (
	"context"
	"errors"

	"lk/internal/models"
	"lk/internal/repository"

	"gorm.io/gorm"
)

// infoService реализует интерфейс InfoService.
type infoService struct {
	serviceRepo repository.ServiceRepository
	infoRepo    repository.InfoRepository
}

// NewInfoService создает новый сервис для получения общей информации.
func NewInfoService(serviceRepo repository.ServiceRepository, infoRepo repository.InfoRepository) InfoService {
	return &infoService{
		serviceRepo: serviceRepo,
		infoRepo:    infoRepo,
	}
}

// GetServiceRecommendations получает рекомендации для конкретной услуги.
func (s *infoService) GetServiceRecommendations(ctx context.Context, serviceID uint64) (models.Recommendation, error) {
	text, err := s.serviceRepo.GetServiceRecommendations(ctx, serviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Recommendation{}, NewNotFoundError(
				"recommendations for this service not found", err)
		}
		return models.Recommendation{}, NewInternalServerError(
			"failed to get service recommendations from db", err)
	}
	return models.Recommendation{Text: text}, nil
}

// GetClinicInfo получает информацию о клинике.
func (s *infoService) GetClinicInfo(ctx context.Context) (models.ClinicInfo, error) {
	info, err := s.infoRepo.GetClinicInfo(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ClinicInfo{}, NewNotFoundError("clinic info not found", err)
		}
		return models.ClinicInfo{}, NewInternalServerError("failed to get clinic info from db", err)
	}
	return info, nil
}

// GetLegalDocuments получает список юридических документов.
func (s *infoService) GetLegalDocuments(ctx context.Context) ([]models.LegalDocument, error) {
	docs, err := s.infoRepo.GetLegalDocuments(ctx)
	if err != nil {
		return nil, NewInternalServerError("failed to get legal documents from db", err)
	}
	return docs, nil
}
