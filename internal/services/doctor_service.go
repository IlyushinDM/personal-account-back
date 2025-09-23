package services

import (
	"context"
	"errors"

	"lk/internal/models"
	"lk/internal/repository"

	"gorm.io/gorm"
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
	doctor, err := s.repo.GetDoctorByID(ctx, doctorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Doctor{}, NewNotFoundError("doctor with this ID not found", err)
		}
		return models.Doctor{}, NewInternalServerError("failed to get doctor details from db", err)
	}
	return doctor, nil
}

// GetSpecialistRecommendations получает рекомендации от врача.
func (s *doctorService) GetSpecialistRecommendations(ctx context.Context, doctorID uint64) (
	models.Recommendation, error,
) {
	text, err := s.repo.GetSpecialistRecommendations(ctx, doctorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Если доктор не найден, возвращаем ошибку 404
			return models.Recommendation{}, NewNotFoundError("specialist with this ID not found", err)
		}
		return models.Recommendation{}, NewInternalServerError("failed to get recommendations from db", err)
	}
	return models.Recommendation{Text: text}, nil
}

// GetDoctorsBySpecialty получает отфильтрованный, отсортированный и пагинированный список врачей.
func (s *doctorService) GetDoctorsBySpecialty(ctx context.Context, specialtyID uint32, params models.PaginationParams) (
	models.PaginatedDoctorsResponse, error,
) {
	doctors, total, err := s.repo.GetDoctorsBySpecialty(ctx, specialtyID, params)
	if err != nil {
		return models.PaginatedDoctorsResponse{}, NewInternalServerError(
			"failed to get doctors by specialty from db", err)
	}
	return models.PaginatedDoctorsResponse{
		Items: doctors,
		Total: total,
	}, nil
}

// SearchDoctors ищет врачей по запросу.
func (s *doctorService) SearchDoctors(ctx context.Context, query string) ([]models.Doctor, error) {
	doctors, err := s.repo.SearchDoctors(ctx, query)
	if err != nil {
		return nil, NewInternalServerError("failed to search doctors by query from db", err)
	}
	return doctors, nil
}

// SearchDoctorsByService ищет врачей по названию услуги.
func (s *doctorService) SearchDoctorsByService(ctx context.Context, serviceQuery string) ([]models.Doctor, error) {
	doctors, err := s.repo.SearchDoctorsByService(ctx, serviceQuery)
	if err != nil {
		return nil, NewInternalServerError("failed to search doctors by service from db", err)
	}
	return doctors, nil
}
