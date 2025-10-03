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
