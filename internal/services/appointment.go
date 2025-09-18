package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"lk/internal/models"
	"lk/internal/repository"
)

var (
	ErrAppointmentNotFound = errors.New("appointment not found")
	ErrForbidden           = errors.New("user does not have permission for this action")
	ErrDoctorNotFound      = errors.New("doctor not found")
)

// appointmentService реализует интерфейс AppointmentService.
type appointmentService struct {
	repo       repository.AppointmentRepository
	doctorRepo repository.DoctorRepository
}

// NewAppointmentService создает новый сервис для управления записями на прием.
func NewAppointmentService(repo repository.AppointmentRepository, doctorRepo repository.DoctorRepository) AppointmentService {
	return &appointmentService{
		repo:       repo,
		doctorRepo: doctorRepo,
	}
}

// CreateAppointment создает новую запись.
func (s *appointmentService) CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error) {
	// 1. Проверить, существует ли такой doctorID.
	_, err := s.doctorRepo.GetDoctorByID(ctx, appointment.DoctorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrDoctorNotFound
		}
		return 0, fmt.Errorf("failed to check doctor existence: %w", err)
	}

	// TODO: Добавить оставшуюся бизнес-логику перед созданием записи:
	// - Проверить, свободен ли врач в это время.
	// - Проверить, существует ли такой serviceID, clinicID.
	// - Отправить уведомление пользователю после успешного создания.
	return s.repo.CreateAppointment(ctx, appointment)
}

// GetUserAppointments получает записи пользователя.
func (s *appointmentService) GetUserAppointments(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	return s.repo.GetAppointmentsByUserID(ctx, userID)
}

// CancelAppointment отменяет запись на прием.
func (s *appointmentService) CancelAppointment(ctx context.Context, userID, appointmentID uint64) error {
	appointment, err := s.repo.GetAppointmentByID(ctx, appointmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAppointmentNotFound
		}
		return err
	}

	// Проверяем, что пользователь отменяет свою собственную запись
	if appointment.UserID != userID {
		return ErrForbidden
	}

	// Например, статус "Отменено пользователем" имеет ID = 3
	const cancelledStatusID = 3
	return s.repo.UpdateAppointmentStatus(ctx, appointmentID, cancelledStatusID)
}

// GetAvailableDates получает доступные для записи даты.
func (s *appointmentService) GetAvailableDates(ctx context.Context, doctorID, serviceID uint64, monthStr string) (models.AvailableDatesResponse, error) {
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return models.AvailableDatesResponse{}, fmt.Errorf("invalid month format, expected YYYY-MM: %w", err)
	}

	// TODO: Добавить логику, использующую serviceID для определения длительности приема и фильтрации дней

	dates, err := s.repo.GetAvailableDates(ctx, doctorID, month)
	if err != nil {
		return models.AvailableDatesResponse{}, err
	}

	stringDates := make([]string, len(dates))
	for i, d := range dates {
		stringDates[i] = d.Format("2006-01-02")
	}

	return models.AvailableDatesResponse{
		SpecialistID:   doctorID,
		Month:          monthStr,
		AvailableDates: stringDates,
	}, nil
}

// GetAvailableSlots получает доступные временные слоты на конкретную дату.
func (s *appointmentService) GetAvailableSlots(ctx context.Context, doctorID, serviceID uint64, dateStr string) (models.AvailableSlotsResponse, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return models.AvailableSlotsResponse{}, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	// TODO: Добавить логику, использующую serviceID для определения длительности приема и фильтрации слотов

	slots, err := s.repo.GetAvailableSlots(ctx, doctorID, date)
	if err != nil {
		return models.AvailableSlotsResponse{}, err
	}

	return models.AvailableSlotsResponse{
		SpecialistID:   doctorID,
		Date:           dateStr,
		AvailableSlots: slots,
	}, nil
}

// GetUpcomingForUser получает предстоящие записи пользователя.
func (s *appointmentService) GetUpcomingForUser(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	return s.repo.GetUpcomingAppointmentsByUserID(ctx, userID)
}
