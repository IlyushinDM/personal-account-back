package services

import (
	"context"
	"errors"

	"lk/internal/models"
	"lk/internal/repository"

	"gorm.io/gorm"
)

// reviewService реализует интерфейс ReviewService.
type reviewService struct {
	repo repository.ReviewRepository
}

// NewReviewService создает новый сервис для работы с отзывами.
func NewReviewService(repo repository.ReviewRepository) ReviewService {
	return &reviewService{repo: repo}
}

// GetReviewsByDoctorID получает пагинированный список отзывов для конкретного врача.
// Возвращает только модерированные отзывы.
func (s *reviewService) GetReviewsByDoctorID(
	ctx context.Context, doctorID uint64, params models.PaginationParams,
) (models.PaginatedReviewsResponse, error) {
	// Валидация параметров пагинации
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10 // Лимит по умолчанию
	}

	// Получаем отзывы из репозитория
	reviews, total, err := s.repo.GetReviewsByDoctorID(ctx, doctorID, params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Если врач не найден или у него нет отзывов, возвращаем пустой результат
			return models.PaginatedReviewsResponse{
				Page:  params.Page,
				Total: 0,
				Items: []models.Review{},
			}, nil
		}
		return models.PaginatedReviewsResponse{}, NewInternalServerError(
			"не удалось получить отзывы из базы данных", err)
	}

	return models.PaginatedReviewsResponse{
		Page:  params.Page,
		Total: total,
		Items: reviews,
	}, nil
}
