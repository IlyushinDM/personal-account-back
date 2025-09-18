package services

import (
	"context"
	"fmt"
	"mime/multipart"

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

// UpdateUserProfile обновляет профиль пользователя.
func (s *userService) UpdateUserProfile(ctx context.Context, userID uint64, input models.UserProfile) (models.UserProfile, error) {
	input.UserID = userID
	return s.userRepo.UpdateUserProfile(ctx, input)
}

// UpdateAvatar обрабатывает логику обновления аватара.
func (s *userService) UpdateAvatar(ctx context.Context, userID uint64, fileHeader *multipart.FileHeader) (string, error) {
	// ! Это mock-реализация.
	// В реальном приложении здесь будет:
	// 1. Генерация уникального имени файла.
	// 2. Загрузка файла в S3-совместимое хранилище (MinIO).
	// 3. Получение URL файла из хранилища.
	// 4. Сохранение URL в базу данных.

	// Имитируем создание URL
	newAvatarURL := "/avatars/" + fmt.Sprintf("%d_%s", userID, fileHeader.Filename)

	err := s.userRepo.UpdateAvatar(ctx, userID, newAvatarURL)
	if err != nil {
		return "", err
	}

	return newAvatarURL, nil
}
