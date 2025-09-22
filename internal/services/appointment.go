package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"lk/internal/models"
	"lk/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrAppointmentNotFound = errors.New("appointment not found")
	ErrForbidden           = errors.New("user does not have permission for this action")
	ErrDoctorNotFound      = errors.New("doctor not found")
	ErrNoSchedule          = errors.New("doctor has no schedule for the selected date")
	ErrNoAvailableSlots    = errors.New("no available slots for the selected date")
)

// appointmentService реализует интерфейс AppointmentService.
type appointmentService struct {
	repo       repository.AppointmentRepository
	doctorRepo repository.DoctorRepository
}

// NewAppointmentService создает новый сервис для управления записями на прием.
func NewAppointmentService(
	repo repository.AppointmentRepository, doctorRepo repository.DoctorRepository,
) AppointmentService {
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
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrDoctorNotFound
		}
		return 0, fmt.Errorf("failed to check doctor existence: %w", err)
	}

	// TODO: Добавить оставшуюся бизнес-логику перед созданием записи:
	// - Проверить, свободен ли врач в это время (самое важное).
	// - Проверить, существует ли такой serviceID, clinicID.
	// - Отправить уведомление пользователю после успешного создания.
	return s.repo.CreateAppointment(ctx, appointment)
}

// GetAvailableDates получает доступные для записи даты.
func (s *appointmentService) GetAvailableDates(ctx context.Context, doctorID, serviceID uint64, monthStr string) (
	models.AvailableDatesResponse, error,
) {
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return models.AvailableDatesResponse{},
			fmt.Errorf("invalid month format, expected YYYY-MM: %w", err)
	}

	// TODO: Добавить логику, использующую serviceID для определения длительности приема и фильтрации дней
	// (например, если в дне не осталось окон под услугу - не показывать его)

	dates, err := s.repo.GetAvailableDatesForMonth(ctx, doctorID, month)
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
func (s *appointmentService) GetAvailableSlots(ctx context.Context, doctorID, serviceID uint64, dateStr string) (
	models.AvailableSlotsResponse, error,
) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return models.AvailableSlotsResponse{},
			fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	// 1. Получить график работы врача на этот день
	schedule, err := s.repo.GetDoctorScheduleForDate(ctx, doctorID, date)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.AvailableSlotsResponse{}, ErrNoSchedule
		}
		return models.AvailableSlotsResponse{}, err
	}

	// 2. Получить длительность запрашиваемой услуги
	requestedServiceDuration, err := s.repo.GetServiceDurationMinutes(ctx, serviceID)
	if err != nil {
		return models.AvailableSlotsResponse{}, fmt.Errorf("could not get service duration: %w", err)
	}
	slotStep := time.Duration(requestedServiceDuration) * time.Minute

	// 3. Получить все уже существующие записи на этот день
	existingAppointments, err := s.repo.GetAppointmentsByDoctorAndDate(ctx, doctorID, date)
	if err != nil {
		return models.AvailableSlotsResponse{}, err
	}

	// 4. Предзагружаем длительности всех услуг для существующих записей (оптимизация)
	var existingServiceIDs []uint64
	if len(existingAppointments) > 0 {
		// Собираем уникальные ID услуг
		serviceIDSet := make(map[uint64]struct{})
		for _, app := range existingAppointments {
			if _, ok := serviceIDSet[app.ServiceID]; !ok {
				serviceIDSet[app.ServiceID] = struct{}{}
				existingServiceIDs = append(existingServiceIDs, app.ServiceID)
			}
		}
	}

	// Делаем один запрос в БД, чтобы получить все нужные длительности
	existingServices, err := s.repo.GetServicesByIDs(ctx, existingServiceIDs)
	if err != nil {
		return models.AvailableSlotsResponse{}, fmt.Errorf("could not get existing services info: %w", err)
	}

	// Создаем мапу для быстрого доступа к длительности
	serviceDurations := make(map[uint64]time.Duration)
	for _, service := range existingServices {
		serviceDurations[service.ID] = time.Duration(service.DurationMinutes) * time.Minute
	}

	// 5. Вычислить доступные слоты
	var availableSlots []string
	loc, _ := time.LoadLocation("UTC") // или другой нужный
	startTime := time.Date(date.Year(), date.Month(), date.Day(), schedule.StartTime.Hour(),
		schedule.StartTime.Minute(), 0, 0, loc)
	endTime := time.Date(date.Year(), date.Month(), date.Day(), schedule.EndTime.Hour(),
		schedule.EndTime.Minute(), 0, 0, loc)

	// Итерируемся по рабочему времени с шагом, равным длительности запрашиваемой услуги
	for slotStart := startTime; slotStart.Add(slotStep).Before(endTime) ||
		slotStart.Add(slotStep).Equal(endTime); slotStart = slotStart.Add(slotStep) {
		slotEnd := slotStart.Add(slotStep)
		isAvailable := true

		// Проверяем, не пересекается ли текущий слот с существующими записями
		for _, app := range existingAppointments {
			appTime, _ := time.Parse("15:04", app.AppointmentTime)
			appStart := time.Date(date.Year(), date.Month(), date.Day(), appTime.Hour(),
				appTime.Minute(), 0, 0, loc)

			// Берем длительность из мапы, а не из БД
			appEnd := appStart.Add(serviceDurations[app.ServiceID])

			// Проверка пересечения интервалов: (StartA < EndB) and (EndA > StartB)
			if slotStart.Before(appEnd) && slotEnd.After(appStart) {
				isAvailable = false
				break
			}
		}

		if isAvailable {
			availableSlots = append(availableSlots, slotStart.Format("15:04"))
		}
	}

	if len(availableSlots) == 0 {
		return models.AvailableSlotsResponse{}, ErrNoAvailableSlots
	}

	return models.AvailableSlotsResponse{
		SpecialistID:   doctorID,
		Date:           dateStr,
		AvailableSlots: availableSlots,
	}, nil
}

func (s *appointmentService) GetUserAppointments(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	return s.repo.GetAppointmentsByUserID(ctx, userID)
}

func (s *appointmentService) CancelAppointment(ctx context.Context, userID, appointmentID uint64) error {
	appointment, err := s.repo.GetAppointmentByID(ctx, appointmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAppointmentNotFound
		}
		return err
	}
	if appointment.UserID != userID {
		return ErrForbidden
	}
	const cancelledStatusID = 3
	return s.repo.UpdateAppointmentStatus(ctx, appointmentID, cancelledStatusID)
}

func (s *appointmentService) GetUpcomingForUser(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	return s.repo.GetUpcomingAppointmentsByUserID(ctx, userID)
}
