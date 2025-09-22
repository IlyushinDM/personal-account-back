package services

import (
	"context"
	"errors"
	"fmt"
	"os"

	"lk/internal/models"
	"lk/internal/repository"

	"gorm.io/gorm"
)

type medicalCardService struct {
	repo             repository.MedicalCardRepository
	prescriptionRepo repository.PrescriptionRepository
}

func NewMedicalCardService(
	repo repository.MedicalCardRepository, prescriptionRepo repository.PrescriptionRepository,
) MedicalCardService {
	return &medicalCardService{
		repo:             repo,
		prescriptionRepo: prescriptionRepo,
	}
}

func (s *medicalCardService) GetVisits(
	ctx context.Context, userID uint64, params models.PaginationParams,
) (models.PaginatedVisitsResponse, error) {
	visits, total, err := s.repo.GetCompletedVisits(ctx, userID, params)
	if err != nil {
		return models.PaginatedVisitsResponse{}, err
	}

	// Обогащаем данные для ответа
	items := make([]models.VisitHistoryItem, 0, len(visits))
	for _, v := range visits {
		item := models.VisitHistoryItem{
			Appointment: v,
			// PatientName - можно взять из UserProfile, если нужно
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

func (s *medicalCardService) GetAnalyses(
	ctx context.Context, userID uint64, status *string,
) ([]models.LabAnalysis, error) {
	return s.repo.GetAnalysesByUserID(ctx, userID, status)
}

func (s *medicalCardService) GetArchivedPrescriptions(
	ctx context.Context, userID uint64,
) ([]models.Prescription, error) {
	return s.repo.GetArchivedPrescriptionsByUserID(ctx, userID)
}

func (s *medicalCardService) GetSummary(
	ctx context.Context, userID uint64,
) (models.MedicalCardSummary, error) {
	return s.repo.GetSummaryInfo(ctx, userID)
}

func (s *medicalCardService) ArchivePrescription(ctx context.Context, userID, prescriptionID uint64) error {
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

	return s.repo.ArchivePrescription(ctx, userID, prescriptionID)
}

// DownloadFile - mock-реализация (FR-4.4)
func (s *medicalCardService) DownloadFile(ctx context.Context, userID, fileID uint64) ([]byte, string, error) {
	// 1. Проверить в БД, что файл с fileID существует и принадлежит userID.
	// Например, найти анализ по fileID и проверить userID.
	// var analysis models.LabAnalysis
	// err := s.repo.GetAnalysisByFileID(fileID) -> if analysis.UserID != userID -> error
	if fileID != 300 { // Mock check
		return nil, "", errors.New("file not found")
	}

	// 2. Получить путь к файлу из БД или S3.
	mockFileName := "Заключение.pdf"
	// ВАЖНО: для работы этой mock-функции нужно создать папку 'static/files' в корне проекта
	// и положить туда файл 'visit_300_report.pdf'
	mockFilePath := "static/files/visit_300_report.pdf"

	// 3. Прочитать файл
	// В реальном приложении здесь будет клиент для S3
	data, err := os.ReadFile(mockFilePath)
	if err != nil {
		return nil, "", fmt.Errorf("could not read mock file: %w", err)
	}

	return data, mockFileName, nil
}
