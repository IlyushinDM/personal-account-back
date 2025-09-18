package services

import (
	"context"

	"lk/internal/models"
	"lk/internal/repository"
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
		return models.Recommendation{}, err
	}
	return models.Recommendation{Text: text}, nil
}

// GetClinicInfo получает информацию о клинике.
func (s *infoService) GetClinicInfo(ctx context.Context) (models.ClinicInfo, error) {
	return s.infoRepo.GetClinicInfo(ctx)
}

// GetLegalDocuments получает список юридических документов.
func (s *infoService) GetLegalDocuments(ctx context.Context) ([]models.LegalDocument, error) {
	return s.infoRepo.GetLegalDocuments(ctx)
}
