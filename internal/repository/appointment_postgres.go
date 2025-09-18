package repository

import (
	"context"
	"time"

	"lk/internal/models"

	"github.com/jmoiron/sqlx"
)

// AppointmentPostgres реализует AppointmentRepository для PostgreSQL.
type AppointmentPostgres struct {
	db *sqlx.DB
}

// NewAppointmentPostgres создает новый экземпляр репозитория для записей на прием.
func NewAppointmentPostgres(db *sqlx.DB) *AppointmentPostgres {
	return &AppointmentPostgres{db: db}
}

// CreateAppointment создает новую запись на прием в базе данных.
func (r *AppointmentPostgres) CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error) {
	var id uint64
	query := `
		INSERT INTO medical_center.appointments (
			user_id, doctor_id, service_id, clinic_id,
			appointment_date, appointment_time, status_id,
			price_at_booking, is_dms
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	row := r.db.QueryRowContext(
		ctx,
		query,
		appointment.UserID,
		appointment.DoctorID,
		appointment.ServiceID,
		appointment.ClinicID,
		appointment.AppointmentDate,
		appointment.AppointmentTime,
		appointment.StatusID,
		appointment.PriceAtBooking,
		appointment.IsDMS,
	)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// GetAppointmentsByUserID получает список всех записей на прием для конкретного пользователя.
func (r *AppointmentPostgres) GetAppointmentsByUserID(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	var appointments []models.Appointment

	query := `
		SELECT
			a.*
		FROM medical_center.appointments a
		WHERE a.user_id=$1
		ORDER BY a.appointment_date DESC, a.appointment_time DESC`

	err := r.db.SelectContext(ctx, &appointments, query, userID)
	return appointments, err
}

// GetAppointmentByID получает запись по ее ID.
func (r *AppointmentPostgres) GetAppointmentByID(ctx context.Context, appointmentID uint64) (models.Appointment, error) {
	var appointment models.Appointment
	query := "SELECT * FROM medical_center.appointments WHERE id=$1"
	err := r.db.GetContext(ctx, &appointment, query, appointmentID)
	return appointment, err
}

// UpdateAppointmentStatus обновляет статус записи.
func (r *AppointmentPostgres) UpdateAppointmentStatus(ctx context.Context, appointmentID uint64, statusID uint32) error {
	query := "UPDATE medical_center.appointments SET status_id = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, statusID, appointmentID)
	return err
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

	query := `
		SELECT a.*
		FROM medical_center.appointments a
		WHERE a.user_id=$1 AND a.appointment_date >= CURRENT_DATE
		ORDER BY a.appointment_date ASC, a.appointment_time ASC`

	err := r.db.SelectContext(ctx, &appointments, query, userID)
	return appointments, err
}
