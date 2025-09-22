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

// NewAppointmentPostgres создает новый экземпляр репозитория для записей на прием.
func NewAppointmentPostgres(db *gorm.DB) *AppointmentPostgres {
	return &AppointmentPostgres{db: db}
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
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("appointment_date desc, appointment_time desc").Find(&appointments).Error
	return appointments, err
}

// GetAppointmentByID получает запись по ее ID.
func (r *AppointmentPostgres) GetAppointmentByID(ctx context.Context, appointmentID uint64) (models.Appointment, error) {
	var appointment models.Appointment
	err := r.db.WithContext(ctx).First(&appointment, appointmentID).Error
	return appointment, err
}

// UpdateAppointmentStatus обновляет статус записи.
func (r *AppointmentPostgres) UpdateAppointmentStatus(ctx context.Context, appointmentID uint64, statusID uint32) error {
	return r.db.WithContext(ctx).Model(&models.Appointment{}).Where("id = ?", appointmentID).Update("status_id", statusID).Error
}

// GetAvailableDates возвращает фиктивный список доступных дат для врача в указанном месяце.
// ! Это mock-реализация. В реальном приложении здесь будет запрос к таблице расписаний.
func (r *AppointmentPostgres) GetAvailableDates(ctx context.Context, doctorID uint64, month time.Time) ([]time.Time, error) {
	// Имитируем, что врач работает 5, 10 и 15 числа каждого месяца
	return []time.Time{
		time.Date(month.Year(), month.Month(), 5, 0, 0, 0, 0, time.UTC),
		time.Date(month.Year(), month.Month(), 10, 0, 0, 0, 0, time.UTC),
		time.Date(month.Year(), month.Month(), 15, 0, 0, 0, 0, time.UTC),
	}, nil
}

// GetAvailableSlots возвращает фиктивный список временных слотов для врача на указанную дату.
// ! Это mock-реализация. В реальном приложении здесь будет сложная логика проверки занятости.
func (r *AppointmentPostgres) GetAvailableSlots(ctx context.Context, doctorID uint64, date time.Time) ([]string, error) {
	return []string{"09:00", "09:30", "11:00", "14:00", "14:30"}, nil
}

// GetUpcomingAppointmentsByUserID получает список предстоящих записей на прием для пользователя.
func (r *AppointmentPostgres) GetUpcomingAppointmentsByUserID(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND appointment_date >= ?", userID, time.Now()).
		Order("appointment_date asc, appointment_time asc").
		Find(&appointments).Error
	return appointments, err
}
