package services

import (
	"context"

	"lk/internal/models"
	"lk/internal/repository"
)

// appointmentService реализует интерфейс AppointmentService.
type appointmentService struct {
	repo repository.AppointmentRepository
}

// NewAppointmentService создает новый сервис для управления записями на прием.
func NewAppointmentService(repo repository.AppointmentRepository) AppointmentService {
	return &appointmentService{repo: repo}
}

// CreateAppointment создает новую запись.
func (s *appointmentService) CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error) {
	// TODO: Добавить бизнес-логику перед созданием записи:
	// 1. Проверить, свободен ли врач в это время.
	// 2. Проверить, существует ли такой serviceID, doctorID, clinicID.
	// 3. Отправить уведомление пользователю после успешного создания.
	return s.repo.CreateAppointment(ctx, appointment)
}

// GetUserAppointments получает записи пользователя.
func (s *appointmentService) GetUserAppointments(ctx context.Context, userID uint64) ([]models.Appointment, error) {
	return s.repo.GetAppointmentsByUserID(ctx, userID)
}
