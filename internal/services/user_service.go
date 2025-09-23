package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"lk/internal/models"
	"lk/internal/repository"
	"lk/internal/storage"

	"github.com/google/uuid"
)

// userService реализует интерфейс UserService.
type userService struct {
	userRepo    repository.UserRepository
	appointRepo repository.AppointmentRepository
	storage     storage.FileStorage
}

// NewUserService создает новый сервис для работы с данными пользователя.
func NewUserService(userRepo repository.UserRepository, appointRepo repository.AppointmentRepository,
	storage storage.FileStorage,
) UserService {
	return &userService{
		userRepo:    userRepo,
		appointRepo: appointRepo,
		storage:     storage,
	}
}

// GetFullUserProfile получает профиль пользователя и его записи на прием.
func (s *userService) GetFullUserProfile(ctx context.Context, userID uint64) (
	models.UserProfile, []models.Appointment, error,
) {
	// Профиль уже загружен в middleware,
	// но логика дублируется на случай вызова сервиса из другого места.
	profile, err := s.userRepo.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return models.UserProfile{}, nil, NewInternalServerError("failed to get user profile from db", err)
	}

	appointments, err := s.appointRepo.GetAppointmentsByUserID(ctx, userID)
	if err != nil {
		return profile, nil, NewInternalServerError("failed to get user appointments from db", err)
	}

	return profile, appointments, nil
}

// UpdateUserProfile обновляет профиль пользователя.
func (s *userService) UpdateUserProfile(ctx context.Context, userID uint64, input models.UserProfile) (
	models.UserProfile, error,
) {
	input.UserID = userID
	updatedProfile, err := s.userRepo.UpdateUserProfile(ctx, input)
	if err != nil {
		return models.UserProfile{}, NewInternalServerError("failed to update user profile in db", err)
	}
	return updatedProfile, nil
}

// UpdateAvatar обрабатывает логику обновления аватара.
func (s *userService) UpdateAvatar(ctx context.Context, userID uint64, fileHeader *multipart.FileHeader) (
	string, error,
) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", NewInternalServerError("failed to open avatar file", err)
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectKey := fmt.Sprintf("avatars/%d/%s%s", userID, uuid.New().String(), ext)

	err = s.storage.Upload(ctx, file, fileHeader.Size, fileHeader.Header.Get(
		"Content-Type"), objectKey)
	if err != nil {
		return "", NewInternalServerError("failed to upload avatar", err)
	}

	err = s.userRepo.UpdateAvatar(ctx, userID, objectKey)
	if err != nil {
		return "", NewInternalServerError("failed to save avatar url to db", err)
	}

	//* Возвращаем полный "URL" или ключ, как решит фронтенд
	// В данном случае, это ключ объекта в MinIO.
	return objectKey, nil
}
