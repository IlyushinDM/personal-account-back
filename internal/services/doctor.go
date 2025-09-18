package services

import (
	"context"

	"lk/internal/models"
	"lk/internal/repository"
)

// doctorService реализует интерфейс DoctorService.
type doctorService struct {
	repo repository.DoctorRepository
}

// NewDoctorService создает новый сервис для работы с врачами.
func NewDoctorService(repo repository.DoctorRepository) DoctorService {
	return &doctorService{repo: repo}
}

// GetDoctorByID получает профиль врача.
func (s *doctorService) GetDoctorByID(ctx context.Context, doctorID uint64) (models.Doctor, error) {
	return s.repo.GetDoctorByID(ctx, doctorID)
}

// GetSpecialistRecommendations получает рекомендации от врача.
func (s *doctorService) GetSpecialistRecommendations(ctx context.Context, doctorID uint64) (models.Recommendation, error) {
	text, err := s.repo.GetSpecialistRecommendations(ctx, doctorID)
	if err != nil {
		return models.Recommendation{}, err
	}
	return models.Recommendation{Text: text}, nil
}

// GetDoctorsBySpecialty получает отфильтрованный, отсортированный и пагинированный список врачей.
func (s *doctorService) GetDoctorsBySpecialty(ctx context.Context, specialtyID uint32, params models.PaginationParams) (models.PaginatedDoctorsResponse, error) {
	doctors, total, err := s.repo.GetDoctorsBySpecialty(ctx, specialtyID, params)
	if err != nil {
		return models.PaginatedDoctorsResponse{}, err
	}
	return models.PaginatedDoctorsResponse{
		Items: doctors,
		Total: total,
	}, nil
}

// SearchDoctors ищет врачей по запросу.
func (s *doctorService) SearchDoctors(ctx context.Context, query string) ([]models.Doctor, error) {
	return s.repo.SearchDoctors(ctx, query)
}

// SearchDoctorsByService ищет врачей по названию услуги.
func (s *doctorService) SearchDoctorsByService(ctx context.Context, serviceQuery string) ([]models.Doctor, error) {
	return s.repo.SearchDoctorsByService(ctx, serviceQuery)
}
