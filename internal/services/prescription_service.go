package services

import (
	"context"
	"database/sql"
	"errors"

	"lk/internal/models"
	"lk/internal/repository"
)

var ErrPrescriptionNotFound = errors.New("prescription not found")

// prescriptionService реализует интерфейс PrescriptionService.
type prescriptionService struct {
	repo repository.PrescriptionRepository
}

// NewPrescriptionService создает новый сервис для работы с назначениями.
func NewPrescriptionService(repo repository.PrescriptionRepository) PrescriptionService {
	return &prescriptionService{repo: repo}
}

// GetActiveForUser получает активные назначения для пользователя.
func (s *prescriptionService) GetActiveForUser(ctx context.Context, userID uint64) ([]models.Prescription, error) {
	return s.repo.GetActiveByUserID(ctx, userID)
}

// ArchiveForUser архивирует назначение для пользователя.
func (s *prescriptionService) ArchiveForUser(ctx context.Context, userID, prescriptionID uint64) error {
	// Проверяем, что назначение существует и принадлежит пользователю
	prescription, err := s.repo.GetByID(ctx, prescriptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPrescriptionNotFound
		}
		return err
	}

	if prescription.UserID != userID {
		return ErrForbidden
	}

	return s.repo.Archive(ctx, userID, prescriptionID)
}
