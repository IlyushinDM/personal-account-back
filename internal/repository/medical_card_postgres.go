package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"lk/internal/models"

	"gorm.io/gorm"
)

type MedicalCardPostgres struct {
	db *gorm.DB
}

func NewMedicalCardPostgres(db *gorm.DB) *MedicalCardPostgres {
	return &MedicalCardPostgres{db: db}
}

// GetCompletedVisits получает завершенные визиты (FR-4.1) с предзагрузкой данных
func (r *MedicalCardPostgres) GetCompletedVisits(ctx context.Context, userID uint64, params models.PaginationParams,
) ([]models.Appointment, int64, error) {
	var visits []models.Appointment
	var total int64
	// Предполагаем, что статус "Завершено" имеет ID = 2
	const completedStatusID = 2

	query := r.db.WithContext(ctx).Model(&models.Appointment{}).
		Where("user_id = ? AND status_id = ?", userID, completedStatusID)

	// Считаем общее количество
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []models.Appointment{}, 0, nil
	}

	// Применяем сортировку и пагинацию
	offset := (params.Page - 1) * params.Limit
	sortOrder := "appointment_date DESC" // Сортировка по умолчанию
	if params.SortBy == "date" && strings.ToUpper(params.SortOrder) == "ASC" {
		sortOrder = "appointment_date ASC"
	}

	// Preload загружает связанные данные одним запросом
	err := query.
		Preload("Doctor").
		Preload("Service").
		Order(sortOrder).
		Limit(params.Limit).
		Offset(offset).
		Find(&visits).Error

	return visits, total, err
}

// GetAnalysesByUserID получает анализы (FR-4.2)
func (r *MedicalCardPostgres) GetAnalysesByUserID(
	ctx context.Context, userID uint64, status *string,
) ([]models.LabAnalysis, error) {
	var analyses []models.LabAnalysis
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if status != nil && *status != "" {
		// Предполагаем, что статусы в БД хранятся как 'completed', 'in_progress'
		query = query.Joins(
			"JOIN medical_center.analysisstatuses ON medical_center.analysisstatuses.id = medical_center.labanalyses.status_id").
			Where("medical_center.analysisstatuses.name = ?", *status)
	}
	err := query.Find(&analyses).Error
	return analyses, err
}

// GetArchivedPrescriptionsByUserID получает архивные назначения (FR-4.3)
func (r *MedicalCardPostgres) GetArchivedPrescriptionsByUserID(
	ctx context.Context, userID uint64,
) ([]models.Prescription, error) {
	var prescriptions []models.Prescription
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, "archived").
		Order("completed_at DESC").
		Find(&prescriptions).Error
	return prescriptions, err
}

// GetSummaryInfo получает сводку (FR-4.6)
func (r *MedicalCardPostgres) GetSummaryInfo(ctx context.Context, userID uint64) (models.MedicalCardSummary, error) {
	var summary models.MedicalCardSummary

	// 1. Последний визит
	var lastVisit struct {
		AppointmentDate time.Time
		DoctorName      string
	}
	err := r.db.WithContext(ctx).Model(&models.Appointment{}).
		Select("appointments.appointment_date, CONCAT(d.last_name, ' ', d.first_name) as doctor_name").
		Joins("JOIN medical_center.doctors d ON d.id = appointments.doctor_id").
		Where("appointments.user_id = ? AND appointments.status_id = ?", userID, 2 /* Завершено */).
		Order("appointments.appointment_date DESC").
		First(&lastVisit).Error

	if err == nil {
		summary.RecentVisit = &models.RecentVisit{
			Date:       lastVisit.AppointmentDate.Format("2006-01-02"),
			DoctorName: lastVisit.DoctorName,
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return summary, err
	}

	// 2. Анализы в работе (предполагаем статус "В работе" имеет ID = 2)
	var pendingAnalyses int64
	err = r.db.WithContext(ctx).Model(&models.LabAnalysis{}).
		Where("user_id = ? AND status_id = ?", userID, 2).
		Count(&pendingAnalyses).Error
	if err != nil {
		return summary, err
	}
	summary.PendingAnalyses = int(pendingAnalyses)

	// 3. Активные назначения
	var activePrescriptions int64
	err = r.db.WithContext(ctx).Model(&models.Prescription{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Count(&activePrescriptions).Error
	if err != nil {
		return summary, err
	}
	summary.ActivePrescriptions = int(activePrescriptions)

	return summary, nil
}

// ArchivePrescription архивирует назначение (FR-4.7)
func (r *MedicalCardPostgres) ArchivePrescription(ctx context.Context, userID, prescriptionID uint64) error {
	result := r.db.WithContext(ctx).Model(&models.Prescription{}).
		Where("id = ? AND user_id = ? AND status = ?", prescriptionID, userID, "active").
		Updates(map[string]interface{}{
			"status":       "archived",
			"completed_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("active prescription not found or does not belong to user")
	}
	return nil
}
