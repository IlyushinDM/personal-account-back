package services

import (
	"context"
	"errors"

	"lk/internal/models"
	"lk/internal/repository"

	"gorm.io/gorm"
)

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
	prescriptions, err := s.repo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, NewInternalServerError("failed to get active prescriptions from db", err)
	}
	return prescriptions, nil
}

// ArchiveForUser архивирует назначение для пользователя.
func (s *prescriptionService) ArchiveForUser(ctx context.Context, userID, prescriptionID uint64) error {
	// Проверяем, что назначение существует
	prescription, err := s.repo.GetByID(ctx, prescriptionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewNotFoundError("prescription not found", err)
		}
		return NewInternalServerError("failed to get prescription from db", err)
	}

	// Проверяем, что оно принадлежит пользователю
	if prescription.UserID != userID {
		return NewForbiddenError("user does not have permission for this action", nil)
	}

	// Выполняем архивацию
	if err := s.repo.Archive(ctx, userID, prescriptionID); err != nil {
		return NewInternalServerError("failed to archive prescription", err)
	}

	return nil
}
