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

// GetReviewByID получает отзыв по ID.
// Возвращает только модерированные отзывы.
func (s *reviewService) GetReviewByID(ctx context.Context, reviewID uint64) (models.Review, error) {
	review, err := s.repo.GetReviewByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Review{}, NewNotFoundError("отзыв не найден", err)
		}
		return models.Review{}, NewInternalServerError(
			"не удалось получить отзыв из базы данных", err)
	}
	return review, nil
}

// CreateReview создает новый отзыв.
func (s *reviewService) CreateReview(ctx context.Context, userID uint64, review models.Review) (uint64, error) {
	// Проверяем, не оставлял ли пользователь уже отзыв этому врачу
	exists, err := s.repo.CheckUserReviewExists(ctx, userID, review.DoctorID)
	if err != nil {
		return 0, NewInternalServerError("не удалось проверить существование отзыва", err)
	}
	if exists {
		return 0, NewConflictError("вы уже оставляли отзыв этому врачу", nil)
	}

	// Устанавливаем пользователя и флаг модерации
	review.UserID = userID
	review.IsModerated = true // Создаем сразу модерированным (временно)

	// Создаем отзыв
	reviewID, err := s.repo.CreateReview(ctx, review)
	if err != nil {
		return 0, NewInternalServerError("не удалось создать отзыв", err)
	}

	return reviewID, nil
}

// UpdateReview обновляет существующий отзыв.
// Пользователь может обновлять только свои отзывы.
func (s *reviewService) UpdateReview(ctx context.Context, userID, reviewID uint64, review models.Review) error {
	// Сначала получаем существующий отзыв для проверки владельца
	existingReview, err := s.repo.GetReviewByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewNotFoundError("отзыв не найден", err)
		}
		return NewInternalServerError("не удалось получить отзыв", err)
	}

	// Проверяем, что отзыв принадлежит пользователю
	if existingReview.UserID != userID {
		return NewForbiddenError("вы можете обновлять только свои отзывы", nil)
	}

	// Обновляем отзыв
	err = s.repo.UpdateReview(ctx, reviewID, review)
	if err != nil {
		return NewInternalServerError("не удалось обновить отзыв", err)
	}

	return nil
}

// DeleteReview удаляет отзыв.
// Пользователь может удалять только свои отзывы.
func (s *reviewService) DeleteReview(ctx context.Context, userID, reviewID uint64) error {
	// Сначала получаем существующий отзыв для проверки владельца
	existingReview, err := s.repo.GetReviewByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewNotFoundError("отзыв не найден", err)
		}
		return NewInternalServerError("не удалось получить отзыв", err)
	}

	// Проверяем, что отзыв принадлежит пользователю
	if existingReview.UserID != userID {
		return NewForbiddenError("вы можете удалять только свои отзывы", nil)
	}

	// Удаляем отзыв
	err = s.repo.DeleteReview(ctx, reviewID)
	if err != nil {
		return NewInternalServerError("не удалось удалить отзыв", err)
	}

	return nil
}

// GetDoctorReviewsWithRating получает отзывы врача с его средним рейтингом.
func (s *reviewService) GetDoctorReviewsWithRating(
	ctx context.Context, doctorID uint64, params models.PaginationParams, onlyModerated bool,
) (models.DoctorReviewsResponse, error) {
	// Валидация параметров пагинации
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	// Получаем отзывы
	reviews, _, err := s.repo.GetReviewsByDoctorID(ctx, doctorID, params)
	if err != nil {
		return models.DoctorReviewsResponse{}, NewInternalServerError(
			"не удалось получить отзывы из базы данных", err)
	}

	// Получаем средний рейтинг врача
	avgRating, err := s.repo.GetDoctorAverageRating(ctx, doctorID)
	if err != nil {
		return models.DoctorReviewsResponse{}, NewInternalServerError(
			"не удалось получить рейтинг врача", err)
	}

	return models.DoctorReviewsResponse{
		Items:        reviews,
		DoctorRating: avgRating,
	}, nil
}
