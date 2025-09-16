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

// GetDoctorsBySpecialty получает список врачей по специальности.
func (s *doctorService) GetDoctorsBySpecialty(ctx context.Context, specialtyID uint32) ([]models.Doctor, error) {
	return s.repo.GetDoctorsBySpecialty(ctx, specialtyID)
}

// SearchDoctors ищет врачей по запросу.
func (s *doctorService) SearchDoctors(ctx context.Context, query string) ([]models.Doctor, error) {
	return s.repo.SearchDoctors(ctx, query)
}
