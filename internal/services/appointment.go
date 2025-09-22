package services

import (
	"context"
	"errors"
	"log"
	"time"

	"lk/internal/models"
	"lk/internal/repository"

	"gorm.io/gorm"
)

// Определяем ошибки как переменные для возможности их проверки через errors.Is
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

	id, err := s.repo.CreateAppointment(ctx, appointment)
	if err != nil {
		return 0, NewInternalServerError("failed to create appointment", err)
	}

	// TODO: FR-5.5 - Инициировать отправку SMS-уведомления.

	return id, nil
}

// GetUserAppointments возвращает все записи пользователя.
func (s *appointmentService) GetUserAppointments(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	appointments, err := s.repo.GetAppointmentsByUserID(ctx, userID)
	if err != nil {
		return nil, NewInternalServerError("failed to get user appointments", err)
	}
	return appointments, nil
}

// CancelAppointment отменяет запись пользователя.
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

	if err := s.repo.UpdateAppointmentStatus(ctx, appointmentID, models.StatusCancelledByPatient); err != nil {
		return NewInternalServerError("failed to update appointment status", err)
	}
	return nil
}

// GetAvailableDates получает доступные для записи даты в месяце.
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

	// 1. Получаем данные из БД для передачи в калькулятор
	schedule, err := s.repo.GetDoctorScheduleForDate(ctx, doctorID, date)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.AvailableSlotsResponse{ // Успешный пустой ответ
				SpecialistID:   doctorID,
				Date:           dateStr,
				AvailableSlots: []string{},
			}, nil
		}
		return models.AvailableSlotsResponse{}, NewInternalServerError("could not get doctor schedule", err)
	}

	existingAppointments, err := s.repo.GetAppointmentsByDoctorAndDate(ctx, doctorID, date)
	if err != nil {
		return models.AvailableSlotsResponse{}, NewInternalServerError("could not get existing appointments", err)
	}

	// 2. Вызываем калькулятор с полученными данными
	slots, err := s.calculateAvailableSlots(ctx, serviceID, date, &schedule, existingAppointments)
	if err != nil {
		if errors.Is(err, ErrNoAvailableSlots) {
			return models.AvailableSlotsResponse{ // Успешный пустой ответ
				SpecialistID:   doctorID,
				Date:           dateStr,
				AvailableSlots: []string{},
			}, nil
		}
		return models.AvailableSlotsResponse{}, err // Пробрасываем другие ошибки (например, Internal)
	}

	return models.AvailableSlotsResponse{
		SpecialistID:   doctorID,
		Date:           dateStr,
		AvailableSlots: slots,
	}, nil
}

// GetAvailableSlotsByRange получает доступные слоты в диапазоне дат.
func (s *appointmentService) GetAvailableSlotsByRange(
	ctx context.Context, doctorID, serviceID uint64, startDateStr, endDateStr string) (
	models.AvailableRangeSlotsResponse, error,
) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return models.AvailableRangeSlotsResponse{}, NewBadRequestError("invalid start date format", err)
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return models.AvailableRangeSlotsResponse{}, NewBadRequestError("invalid end date format", err)
	}
	if startDate.After(endDate) {
		return models.AvailableRangeSlotsResponse{}, NewBadRequestError("start date cannot be after end date", nil)
	}

	schedules, err := s.repo.GetDoctorScheduleForDateRange(ctx, doctorID, startDate, endDate)
	if err != nil {
		return models.AvailableRangeSlotsResponse{}, NewInternalServerError("could not get schedules for range", err)
	}

	existingAppointments, err := s.repo.GetAppointmentsByDoctorAndDateRange(ctx, doctorID, startDate, endDate)
	if err != nil {
		return models.AvailableRangeSlotsResponse{}, NewInternalServerError("could not get appointments for range", err)
	}

	appointmentsByDate := make(map[time.Time][]models.Appointment)
	for _, app := range existingAppointments {
		day := time.Date(app.AppointmentDate.Year(), app.AppointmentDate.Month(),
			app.AppointmentDate.Day(), 0, 0, 0, 0, time.UTC)
		appointmentsByDate[day] = append(appointmentsByDate[day], app)
	}

	var slotsByDay []models.SlotsForDay
	for _, schedule := range schedules {
		day := time.Date(schedule.Date.Year(), schedule.Date.Month(), schedule.Date.Day(),
			0, 0, 0, 0, time.UTC)
		slots, err := s.calculateAvailableSlots(
			ctx, serviceID, day, &schedule, appointmentsByDate[day])
		if err != nil && !errors.Is(err, ErrNoAvailableSlots) {
			log.Printf(
				"WARN: could not calculate slots for date %s: %v", day.Format("2006-01-02"), err)
			continue
		}

		if len(slots) > 0 {
			slotsByDay = append(slotsByDay, models.SlotsForDay{
				Date:           day.Format("2006-01-02"),
				AvailableSlots: slots,
			})
		}
	}

	return models.AvailableRangeSlotsResponse{
		SpecialistID: doctorID,
		ServiceID:    serviceID,
		SlotsByDay:   slotsByDay,
	}, nil
}

// GetUpcomingForUser получает предстоящие записи пользователя.
func (s *appointmentService) GetUpcomingForUser(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	appointments, err := s.repo.GetUpcomingAppointmentsByUserID(ctx, userID)
	if err != nil {
		return nil, NewInternalServerError("failed to get upcoming appointments", err)
	}
	return appointments, nil
}

// calculateAvailableSlots инкапсулирует логику расчета слотов.
func (s *appointmentService) calculateAvailableSlots(
	ctx context.Context, serviceID uint64, date time.Time,
	schedule *models.Schedule, existingAppointments []models.Appointment,
) ([]string, error) {
	if schedule == nil {
		return nil, ErrNoSchedule
	}

	requestedServiceDuration, err := s.repo.GetServiceDurationMinutes(ctx, serviceID)
	if err != nil {
		return nil, NewInternalServerError("could not get service duration", err)
	}
	slotStep := time.Duration(requestedServiceDuration) * time.Minute

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
