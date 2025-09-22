package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"lk/internal/models"
	"lk/internal/repository"
	"lk/internal/storage"

	"gorm.io/gorm"
)

// medicalCardService реализует интерфейс MedicalCardService.
type medicalCardService struct {
	repo             repository.MedicalCardRepository
	prescriptionRepo repository.PrescriptionRepository
	storage          storage.FileStorage
}

// NewMedicalCardService создает новый сервис для работы с медкартой.
func NewMedicalCardService(
	repo repository.MedicalCardRepository,
	prescriptionRepo repository.PrescriptionRepository,
	storage storage.FileStorage,
) MedicalCardService {
	return &medicalCardService{
		repo:             repo,
		prescriptionRepo: prescriptionRepo,
		storage:          storage,
	}
}

// GetVisits получает историю завершенных визитов пользователя.
func (s *medicalCardService) GetVisits(ctx context.Context, userID uint64, params models.PaginationParams) (models.PaginatedVisitsResponse, error) {
	visits, total, err := s.repo.GetCompletedVisits(ctx, userID, params)
	if err != nil {
		return models.PaginatedVisitsResponse{}, NewInternalServerError("failed to get visits from db", err)
	}

	items := make([]models.VisitHistoryItem, 0, len(visits))
	for _, v := range visits {
		item := models.VisitHistoryItem{
			Appointment: v,
			DoctorName:  fmt.Sprintf("%s %s", v.Doctor.LastName, v.Doctor.FirstName),
			ServiceName: v.Service.Name,
		}
		items = append(items, item)
	}

	return models.PaginatedVisitsResponse{
		Total: total,
		Items: items,
	}, nil
}

// GetAnalyses получает список анализов пользователя.
func (s *medicalCardService) GetAnalyses(ctx context.Context, userID uint64, status *string) ([]models.LabAnalysis, error) {
	analyses, err := s.repo.GetAnalysesByUserID(ctx, userID, status)
	if err != nil {
		return nil, NewInternalServerError("failed to get analyses from db", err)
	}
	return analyses, nil
}

// GetArchivedPrescriptions получает архивные назначения пользователя.
func (s *medicalCardService) GetArchivedPrescriptions(ctx context.Context, userID uint64) ([]models.Prescription, error) {
	prescriptions, err := s.repo.GetArchivedPrescriptionsByUserID(ctx, userID)
	if err != nil {
		return nil, NewInternalServerError("failed to get archived prescriptions from db", err)
	}
	return prescriptions, nil
}

// GetSummary получает сводную информацию по медкарте.
func (s *medicalCardService) GetSummary(ctx context.Context, userID uint64) (models.MedicalCardSummary, error) {
	summary, err := s.repo.GetSummaryInfo(ctx, userID)
	if err != nil {
		return models.MedicalCardSummary{}, NewInternalServerError("failed to get summary info from db", err)
	}
	return summary, nil
}

// ArchivePrescription архивирует назначение, проверяя его принадлежность пользователю.
func (s *medicalCardService) ArchivePrescription(ctx context.Context, userID, prescriptionID uint64) error {
	prescription, err := s.prescriptionRepo.GetByID(ctx, prescriptionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewNotFoundError("prescription not found", err)
		}
		return NewInternalServerError("failed to get prescription from db", err)
	}

	if prescription.UserID != userID {
		return NewForbiddenError("user does not have permission for this action", nil)
	}

	if err := s.repo.ArchivePrescription(ctx, userID, prescriptionID); err != nil {
		return NewInternalServerError("failed to archive prescription", err)
	}
	return nil
}

// DownloadFile находит запись о файле в БД, проверяет права доступа и возвращает содержимое файла.
func (s *medicalCardService) DownloadFile(ctx context.Context, userID, analysisID uint64) ([]byte, string, error) {
	analysis, err := s.repo.GetAnalysisByID(ctx, analysisID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", NewNotFoundError("file record not found", err)
		}
		return nil, "", NewInternalServerError("failed to get analysis from db", err)
	}

	if analysis.UserID != userID {
		return nil, "", NewForbiddenError("user does not have permission to access this file", nil)
	}

	if !analysis.ResultFileURL.Valid || analysis.ResultFileURL.String == "" {
		return nil, "", NewNotFoundError("file path is not specified for this analysis", nil)
	}

	objectKey := analysis.ResultFileURL.String
	fileName := analysis.ResultFileName.String
	if fileName == "" {
		fileName = "download" // Имя по умолчанию
	}

	fileObject, err := s.storage.Download(ctx, objectKey)
	if err != nil {
		return nil, "", NewInternalServerError("could not get file from storage", err)
	}
	defer fileObject.Close()

	data, err := io.ReadAll(fileObject)
	if err != nil {
		log.Printf("ERROR: Could not read file stream for key %s: %v", objectKey, err)
		return nil, "", NewInternalServerError("could not read file stream", err)
	}

	return data, fileName, nil
}
