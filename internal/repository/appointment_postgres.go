package repository

import (
	"context"
	"time"

	"lk/internal/models"

	"gorm.io/gorm"
)

// AppointmentPostgres реализует AppointmentRepository для PostgreSQL.
type AppointmentPostgres struct {
	db *gorm.DB
}

// NewAppointmentPostgres создает новый экземпляр репозитория.
func NewAppointmentPostgres(db *gorm.DB) *AppointmentPostgres {
	return &AppointmentPostgres{db: db}
}

// GetDoctorScheduleForDate получает расписание врача на конкретную дату.
func (r *AppointmentPostgres) GetDoctorScheduleForDate(
	ctx context.Context, doctorID uint64, date time.Time,
) (models.Schedule, error) {
	var schedule models.Schedule
	err := r.db.WithContext(ctx).Where(
		"doctor_id = ? AND date = ?", doctorID, date).First(&schedule).Error
	return schedule, err
}

// GetAppointmentsByDoctorAndDate получает все записи к врачу на конкретную дату.
func (r *AppointmentPostgres) GetAppointmentsByDoctorAndDate(
	ctx context.Context, doctorID uint64, date time.Time,
) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.WithContext(ctx).Where(
		"doctor_id = ? AND appointment_date = ?", doctorID, date).Find(&appointments).Error
	return appointments, err
}

// GetAppointmentsByDoctorAndDateRange получает все записи к врачу в диапазоне дат.
func (r *AppointmentPostgres) GetAppointmentsByDoctorAndDateRange(
	ctx context.Context, doctorID uint64, startDate, endDate time.Time,
) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.WithContext(ctx).Where(
		"doctor_id = ? AND appointment_date BETWEEN ? AND ?",
		doctorID, startDate, endDate,
	).Find(&appointments).Error
	return appointments, err
}

// GetDoctorScheduleForDateRange получает все расписания врача в диапазоне дат.
func (r *AppointmentPostgres) GetDoctorScheduleForDateRange(
	ctx context.Context, doctorID uint64, startDate, endDate time.Time,
) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.WithContext(ctx).Where(
		"doctor_id = ? AND date BETWEEN ? AND ?",
		doctorID, startDate, endDate,
	).Order("date ASC").Find(&schedules).Error
	return schedules, err
}

// GetServiceDurationMinutes получает длительность услуги в минутах.
func (r *AppointmentPostgres) GetServiceDurationMinutes(ctx context.Context, serviceID uint64) (uint16, error) {
	var service models.Service
	err := r.db.WithContext(ctx).Select("duration_minutes").First(&service, serviceID).Error
	return service.DurationMinutes, err
}

// GetServicesByIDs получает информацию о нескольких услугах одним запросом.
func (r *AppointmentPostgres) GetServicesByIDs(ctx context.Context, serviceIDs []uint64) ([]models.Service, error) {
	var services []models.Service
	if len(serviceIDs) == 0 {
		return services, nil // Возвращаем пустой слайс, если нет ID для поиска
	}
	err := r.db.WithContext(ctx).Where("id IN ?", serviceIDs).Find(&services).Error
	return services, err
}

// GetAvailableDatesForMonth возвращает дни, в которые у врача есть расписание.
func (r *AppointmentPostgres) GetAvailableDatesForMonth(
	ctx context.Context, doctorID uint64, month time.Time,
) ([]time.Time, error) {
	var dates []time.Time
	startOfMonth := month.Format("2006-01-02")
	endOfMonth := month.AddDate(0, 1, -1).Format("2006-01-02")

	err := r.db.WithContext(ctx).Model(&models.Schedule{}).
		Distinct("date").
		Where("doctor_id = ? AND date BETWEEN ? AND ?", doctorID, startOfMonth, endOfMonth).
		Order("date").
		Pluck("date", &dates).Error

	return dates, err
}

// CreateAppointment создает новую запись на прием в базе данных.
func (r *AppointmentPostgres) CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error) {
	result := r.db.WithContext(ctx).Create(&appointment)
	if result.Error != nil {
		return 0, result.Error
	}
	return appointment.ID, nil
}

// GetAppointmentsByUserID получает список всех записей на прием для конкретного пользователя.
func (r *AppointmentPostgres) GetAppointmentsByUserID(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.WithContext(ctx).Where(
		"user_id = ?", userID).Order("appointment_date desc, appointment_time desc").Find(
		&appointments).Error
	return appointments, err
}

// GetAppointmentByID получает запись по ее ID.
func (r *AppointmentPostgres) GetAppointmentByID(ctx context.Context, appointmentID uint64) (
	models.Appointment, error,
) {
	var appointment models.Appointment
	err := r.db.WithContext(ctx).First(&appointment, appointmentID).Error
	return appointment, err
}

// UpdateAppointmentStatus обновляет статус записи.
func (r *AppointmentPostgres) UpdateAppointmentStatus(ctx context.Context, appointmentID uint64, statusID uint32) error {
	return r.db.WithContext(ctx).Model(&models.Appointment{}).Where(
		"id = ?", appointmentID).Update("status_id", statusID).Error
}

// GetUpcomingAppointmentsByUserID получает список предстоящих записей на прием для пользователя.
func (r *AppointmentPostgres) GetUpcomingAppointmentsByUserID(ctx context.Context, userID uint64) (
	[]models.Appointment, error,
) {
	var appointments []models.Appointment
	// Сравниваем только дату
	today := time.Now().Format("2006-01-02")
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND appointment_date >= ?", userID, today).
		Order("appointment_date asc, appointment_time asc").
		Find(&appointments).Error
	return appointments, err
}
