package services

import (
	"context"
	"errors"
	"time"

	"lk/internal/models"
	"lk/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrNoSchedule       = errors.New("doctor has no schedule for the selected date")
	ErrNoAvailableSlots = errors.New("no available slots for the selected date")
)

// appointmentService реализует интерфейс AppointmentService.
type appointmentService struct {
	repo       repository.AppointmentRepository
	doctorRepo repository.DoctorRepository
	location   *time.Location
}

// NewAppointmentService создает новый сервис для управления записями на прием.
func NewAppointmentService(
	repo repository.AppointmentRepository,
	doctorRepo repository.DoctorRepository,
	location *time.Location,
) AppointmentService {
	return &appointmentService{
		repo:       repo,
		doctorRepo: doctorRepo,
		location:   location,
	}
}

// CreateAppointment создает новую запись.
func (s *appointmentService) CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error) {
	// 1. Проверить, существует ли такой doctorID.
	_, err := s.doctorRepo.GetDoctorByID(ctx, appointment.DoctorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, NewNotFoundError("doctor not found", err)
		}
		return 0, NewInternalServerError("failed to check doctor existence", err)
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
			NewBadRequestError("invalid month format, expected YYYY-MM", err)
	}

	dates, err := s.repo.GetAvailableDatesForMonth(ctx, doctorID, month)
	if err != nil {
		return models.AvailableDatesResponse{}, NewInternalServerError("failed to get available dates", err)
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
			NewBadRequestError("invalid date format, expected YYYY-MM-DD", err)
	}

	// Выносим основную логику в отдельную функцию для чистоты
	slots, err := s.calculateAvailableSlots(ctx, doctorID, serviceID, date)
	if err != nil {
		if errors.Is(err, ErrNoSchedule) || errors.Is(err, ErrNoAvailableSlots) {
			return models.AvailableSlotsResponse{
				SpecialistID:   doctorID,
				Date:           dateStr,
				AvailableSlots: []string{},
			}, nil
		}
		return models.AvailableSlotsResponse{}, err
	}

	return models.AvailableSlotsResponse{
		SpecialistID:   doctorID,
		Date:           dateStr,
		AvailableSlots: slots,
	}, nil
}

func (s *appointmentService) GetUserAppointments(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	return s.repo.GetAppointmentsByUserID(ctx, userID)
}

func (s *appointmentService) CancelAppointment(ctx context.Context, userID, appointmentID uint64) error {
	appointment, err := s.repo.GetAppointmentByID(ctx, appointmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewNotFoundError("appointment not found", err)
		}
		return NewInternalServerError("failed to get appointment", err)
	}

	if appointment.UserID != userID {
		return NewForbiddenError("user does not have permission for this action", nil)
	}

	const cancelledStatusID = 3
	if err := s.repo.UpdateAppointmentStatus(ctx, appointmentID, cancelledStatusID); err != nil {
		return NewInternalServerError("failed to update appointment status", err)
	}
	return nil
}

func (s *appointmentService) GetUpcomingForUser(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	return s.repo.GetUpcomingAppointmentsByUserID(ctx, userID)
}

// calculateAvailableSlots инкапсулирует логику расчета слотов.
// Возвращает слайс строк или одну из кастомных ошибок.
func (s *appointmentService) calculateAvailableSlots(
	ctx context.Context, doctorID, serviceID uint64, date time.Time,
) ([]string, error) {
	// 1. Получить график работы врача на этот день
	schedule, err := s.repo.GetDoctorScheduleForDate(ctx, doctorID, date)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoSchedule
		}
		return nil, NewInternalServerError("could not get doctor schedule", err)
	}

	// 2. Получить длительность запрашиваемой услуги
	requestedServiceDuration, err := s.repo.GetServiceDurationMinutes(ctx, serviceID)
	if err != nil {
		return nil, NewInternalServerError("could not get service duration", err)
	}
	slotStep := time.Duration(requestedServiceDuration) * time.Minute

	// 3. Получить все уже существующие записи на этот день
	existingAppointments, err := s.repo.GetAppointmentsByDoctorAndDate(ctx, doctorID, date)
	if err != nil {
		return nil, NewInternalServerError("could not get existing appointments", err)
	}

	// 4. Предзагружаем длительности всех услуг для существующих записей (оптимизация)
	var existingServiceIDs []uint64
	if len(existingAppointments) > 0 {
		serviceIDSet := make(map[uint64]struct{})
		for _, app := range existingAppointments {
			if _, ok := serviceIDSet[app.ServiceID]; !ok {
				serviceIDSet[app.ServiceID] = struct{}{}
				existingServiceIDs = append(existingServiceIDs, app.ServiceID)
			}
		}
	}

	existingServices, err := s.repo.GetServicesByIDs(ctx, existingServiceIDs)
	if err != nil {
		return nil, NewInternalServerError("could not get existing services info", err)
	}

	serviceDurations := make(map[uint64]time.Duration)
	for _, service := range existingServices {
		serviceDurations[service.ID] = time.Duration(service.DurationMinutes) * time.Minute
	}

	// 5. Вычислить доступные слоты
	var availableSlots []string
	startTime := time.Date(date.Year(), date.Month(), date.Day(), schedule.StartTime.Hour(),
		schedule.StartTime.Minute(), 0, 0, s.location)
	endTime := time.Date(date.Year(), date.Month(), date.Day(), schedule.EndTime.Hour(),
		schedule.EndTime.Minute(), 0, 0, s.location)

	for slotStart := startTime; slotStart.Add(slotStep).Before(endTime) ||
		slotStart.Add(slotStep).Equal(endTime); slotStart = slotStart.Add(slotStep) {
		slotEnd := slotStart.Add(slotStep)
		isAvailable := true

		for _, app := range existingAppointments {
			appTime, _ := time.Parse("15:04", app.AppointmentTime)
			appStart := time.Date(date.Year(), date.Month(), date.Day(), appTime.Hour(),
				appTime.Minute(), 0, 0, s.location)

			appEnd := appStart.Add(serviceDurations[app.ServiceID])

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
		return nil, ErrNoAvailableSlots
	}

	return availableSlots, nil
}
