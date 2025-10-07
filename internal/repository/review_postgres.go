package repository

import (
	"context"
	"fmt"
	"strings"

	"lk/internal/models"

	"gorm.io/gorm"
)

// ReviewPostgres реализует интерфейс ReviewRepository для PostgreSQL.
type ReviewPostgres struct {
	db *gorm.DB
}

// NewReviewPostgres создает новый экземпляр репозитория для отзывов.
func NewReviewPostgres(db *gorm.DB) *ReviewPostgres {
	return &ReviewPostgres{db: db}
}

// GetReviewsByDoctorID получает список отзывов по ID врача с пагинацией и сортировкой.
// Возвращает только модерированные отзывы (is_moderated = true).
func (r *ReviewPostgres) GetReviewsByDoctorID(
	ctx context.Context, doctorID uint64, params models.PaginationParams,
) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	// мапа для масштабируемой сортировки, по хорошему нужно еще добавить updated_at
	allowedSortBy := map[string]string{
		"rating":     "rating",
		"created_at": "created_at",
		"date":       "created_at", // альтернативное название для created_at
	}

	// Получаем колонку для сортировки, по умолчанию - rating
	orderByColumn, ok := allowedSortBy[params.SortBy]
	if !ok {
		orderByColumn = "rating" // Сортировка по умолчанию
	}

	// Определяем порядок сортировки
	sortOrder := "DESC"
	if strings.ToUpper(params.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	// Базовый запрос с фильтрацией по врачу и модерированным отзывам
	query := r.db.WithContext(ctx).Model(&models.Review{}).
		Where("doctor_id = ? AND is_moderated = ?", doctorID, true)

	// Получаем общее количество модерированных отзывов для данного врача
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Если отзывов нет, возвращаем пустой результат
	if total == 0 {
		return []models.Review{}, 0, nil
	}

	// Получаем пагинированные данные
	offset := (params.Page - 1) * params.Limit
	orderClause := fmt.Sprintf("%s %s", orderByColumn, sortOrder)

	err := query.Order(orderClause).
		Limit(params.Limit).
		Offset(offset).
		Find(&reviews).Error

	return reviews, total, err
}

// GetReviewByID получает отзыв по ID.
// Возвращает только модерированные отзывы.
func (r *ReviewPostgres) GetReviewByID(ctx context.Context, reviewID uint64) (models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_moderated = ?", reviewID, true).
		First(&review).Error
	return review, err
}

// CreateReview создает новый отзыв.
func (r *ReviewPostgres) CreateReview(ctx context.Context, review models.Review) (uint64, error) {
	err := r.db.WithContext(ctx).Create(&review).Error
	return review.ID, err
}

// UpdateReview обновляет существующий отзыв.
func (r *ReviewPostgres) UpdateReview(ctx context.Context, reviewID uint64, review models.Review) error {
	return r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("id = ?", reviewID).
		Updates(map[string]interface{}{
			"rating":       review.Rating,
			"comment":      review.Comment,
			"is_moderated": false, // Сбрасываем флаг модерации при обновлении
		}).Error
}

// DeleteReview удаляет отзыв по ID.
func (r *ReviewPostgres) DeleteReview(ctx context.Context, reviewID uint64) error {
	return r.db.WithContext(ctx).
		Where("id = ?", reviewID).
		Delete(&models.Review{}).Error
}

// CheckUserReviewExists проверяет, существует ли уже отзыв от пользователя для данного врача.
func (r *ReviewPostgres) CheckUserReviewExists(ctx context.Context, userID, doctorID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("user_id = ? AND doctor_id = ?", userID, doctorID).
		Count(&count).Error
	return count > 0, err
}

// GetDoctorAverageRating получает средний рейтинг врача из всех отзывов.
func (r *ReviewPostgres) GetDoctorAverageRating(ctx context.Context, doctorID uint64) (float64, error) {
	var avgRating float64
	err := r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("doctor_id = ? AND is_moderated = ?", doctorID, true).
		Select("AVG(rating)").
		Scan(&avgRating).Error
	return avgRating, err
}
