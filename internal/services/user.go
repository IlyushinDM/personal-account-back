package services

import (
	"context"

	"lk/internal/models"
	"lk/internal/repository"
)

// userService реализует интерфейс UserService.
type userService struct {
	userRepo    repository.UserRepository
	appointRepo repository.AppointmentRepository
}

// NewUserService создает новый сервис для работы с данными пользователя.
func NewUserService(userRepo repository.UserRepository, appointRepo repository.AppointmentRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		appointRepo: appointRepo,
	}
}

// GetFullUserProfile получает профиль пользователя и его записи на прием.
func (s *userService) GetFullUserProfile(ctx context.Context, userID uint64) (
	models.UserProfile, []models.Appointment, error,
) {
	profile, err := s.userRepo.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return models.UserProfile{}, nil, err
	}

	appointments, err := s.appointRepo.GetAppointmentsByUserID(ctx, userID)
	if err != nil {
		return profile, nil, err
	}

	return profile, appointments, nil
}
