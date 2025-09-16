package repository

import (
	"context"

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
