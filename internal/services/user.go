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
func NewUserService(userRepo repository.UserRepository,
	appointRepo repository.AppointmentRepository, storage storage.FileStorage,
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
func (s *userService) UpdateUserProfile(
	ctx context.Context, userID uint64, input models.UserProfile,
) (models.UserProfile, error) {
	input.UserID = userID
	return s.userRepo.UpdateUserProfile(ctx, input)
}

// UpdateAvatar обрабатывает логику обновления аватара.
func (s *userService) UpdateAvatar(
	ctx context.Context, userID uint64, fileHeader *multipart.FileHeader,
) (string, error) {
	// 1. Открываем файл для чтения
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open avatar file: %w", err)
	}
	defer file.Close()

	// 2. Генерируем уникальное имя файла (ключ объекта).
	ext := filepath.Ext(fileHeader.Filename)
	objectKey := fmt.Sprintf("avatars/%d/%s%s", userID, uuid.New().String(), ext)

	// 3. Загружаем файл в S3-хранилище (MinIO).
	err = s.storage.Upload(
		ctx, file, fileHeader.Size, fileHeader.Header.Get("Content-Type"), objectKey)
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}

	// 4. Сохраняем ключ объекта в базу данных.
	// * Примечание: Мы сохраняем только ключ, а не полный URL. URL может быть сформирован на клиенте или здесь.
	err = s.userRepo.UpdateAvatar(ctx, userID, objectKey)
	if err != nil {
		// TODO: Здесь нужна логика удаления файла из MinIO, если не удалось обновить БД.
		return "", fmt.Errorf("failed to save avatar url to db: %w", err)
	}

	// Возвращаем ключ объекта (или можно сформировать полный URL, если нужно).
	return objectKey, nil
}
