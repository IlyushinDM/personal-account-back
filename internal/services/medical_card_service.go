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

type medicalCardService struct {
	repo             repository.MedicalCardRepository
	prescriptionRepo repository.PrescriptionRepository
	storage          storage.FileStorage
}

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
		return models.PaginatedVisitsResponse{}, err
	}

	// Обогащаем данные для ответа DTO
	items := make([]models.VisitHistoryItem, 0, len(visits))
	for _, v := range visits {
		item := models.VisitHistoryItem{
			Appointment: v,
			// PatientName можно будет добавить, если потребуется получать профиль пользователя
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
	return s.repo.GetAnalysesByUserID(ctx, userID, status)
}

// GetArchivedPrescriptions получает архивные назначения пользователя.
func (s *medicalCardService) GetArchivedPrescriptions(ctx context.Context, userID uint64) ([]models.Prescription, error) {
	return s.repo.GetArchivedPrescriptionsByUserID(ctx, userID)
}

// GetSummary получает сводную информацию по медкарте.
func (s *medicalCardService) GetSummary(ctx context.Context, userID uint64) (models.MedicalCardSummary, error) {
	return s.repo.GetSummaryInfo(ctx, userID)
}

// ArchivePrescription архивирует назначение, проверяя его принадлежность пользователю.
func (s *medicalCardService) ArchivePrescription(ctx context.Context, userID, prescriptionID uint64) error {
	// 1. Проверяем, что назначение существует и принадлежит пользователю
	prescription, err := s.prescriptionRepo.GetByID(ctx, prescriptionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPrescriptionNotFound
		}
		return err
	}

	if prescription.UserID != userID {
		return ErrForbidden
	}

	// 2. Выполняем архивацию
	return s.repo.ArchivePrescription(ctx, userID, prescriptionID)
}

// DownloadFile находит запись о файле в БД, проверяет права и скачивает из S3.
// Возвращает поток данных, имя файла и ошибку.
func (s *medicalCardService) DownloadFile(ctx context.Context, userID, analysisID uint64) ([]byte, string, error) {
	analysis, err := s.repo.GetAnalysisByID(ctx, analysisID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("file record not found")
		}
		return nil, "", err
	}
	if analysis.UserID != userID {
		return nil, "", ErrForbidden
	}
	if !analysis.ResultFileURL.Valid || analysis.ResultFileURL.String == "" {
		return nil, "", errors.New("file path is not specified for this analysis")
	}

	objectKey := analysis.ResultFileURL.String
	fileName := analysis.ResultFileName.String
	if fileName == "" {
		fileName = "download" // Имя по умолчанию
	}

	// Скачиваем файл из MinIO
	fileObject, err := s.storage.Download(ctx, objectKey)
	if err != nil {
		log.Printf("ERROR: Could not download file from storage with key %s: %v", objectKey, err)
		return nil, "", fmt.Errorf("could not get file from storage")
	}
	defer fileObject.Close()

	// Читаем все байты из потока
	data, err := io.ReadAll(fileObject)
	if err != nil {
		log.Printf("ERROR: Could not read file stream for key %s: %v", objectKey, err)
		return nil, "", fmt.Errorf("could not read file stream")
	}

	return data, fileName, nil
}
