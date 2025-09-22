package repository

import (
	"context"
	"fmt"
	"strings"

	"lk/internal/models"

	"gorm.io/gorm"
)

// DoctorPostgres реализует интерфейс DoctorRepository для PostgreSQL.
type DoctorPostgres struct {
	db *gorm.DB
}

// NewDoctorPostgres создает новый экземпляр репозитория для врачей.
func NewDoctorPostgres(db *gorm.DB) *DoctorPostgres {
	return &DoctorPostgres{db: db}
}

// GetDoctorByID получает информацию о враче по его ID.
func (r *DoctorPostgres) GetDoctorByID(ctx context.Context, id uint64) (models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.WithContext(ctx).Preload("Specialty").First(&doctor, id).Error
	return doctor, err
}

// GetSpecialistRecommendations получает рекомендации из поля в таблице doctors.
func (r *DoctorPostgres) GetSpecialistRecommendations(ctx context.Context, doctorID uint64) (string, error) {
	var doctor models.Doctor
	// Выбираем только одно поле для эффективности
	err := r.db.WithContext(ctx).Select("recommendations").First(&doctor, doctorID).Error
	if err != nil {
		return "", err
	}
	return doctor.Recommendations.String, nil
}

// GetDoctorsBySpecialty получает список врачей по ID их специальности с пагинацией и сортировкой.
func (r *DoctorPostgres) GetDoctorsBySpecialty(
	ctx context.Context, specialtyID uint32, params models.PaginationParams,
) ([]models.Doctor, int64, error) {
	var doctors []models.Doctor
	var total int64

	// Белый список для колонок сортировки для предотвращения SQL-инъекций
	allowedSortBy := map[string]string{
		"rating":           "rating",
		"experience":       "experience_years",
		"experience_years": "experience_years",
		"name":             "last_name",
	}
	orderByColumn, ok := allowedSortBy[params.SortBy]
	if !ok {
		orderByColumn = "rating" // Сортировка по умолчанию
	}

	sortOrder := "DESC"
	if strings.ToUpper(params.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	query := r.db.WithContext(ctx).Model(&models.Doctor{}).Where("specialty_id = ?", specialtyID)

	// Получаем общее количество
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []models.Doctor{}, 0, nil
	}

	// Получаем пагинированные данные
	offset := (params.Page - 1) * params.Limit
	orderClause := fmt.Sprintf("%s %s", orderByColumn, sortOrder)
	err := query.Preload("Specialty").Order(
		orderClause).Limit(params.Limit).Offset(offset).Find(&doctors).Error

	return doctors, total, err
}

// SearchDoctors выполняет поиск врачей по ФИО или названию специальности с использованием FTS.
func (r *DoctorPostgres) SearchDoctors(ctx context.Context, searchQuery string) ([]models.Doctor, error) {
	var doctors []models.Doctor

	// Преобразуем поисковый запрос в формат, понятный plainto_tsquery
	// Заменяем пробелы на оператор 'И' (&) для поиска всех слов
	tsQuery := strings.ReplaceAll(strings.TrimSpace(searchQuery), " ", " & ")

	// Ранжируем результаты по релевантности
	rank := "ts_rank(fts_document, plainto_tsquery('russian', ?))"

	err := r.db.WithContext(ctx).
		Preload("Specialty").
		Select("*, "+rank+" as rank", searchQuery).
		Where("fts_document @@ plainto_tsquery('russian', ?)", tsQuery).
		Order("rank DESC").
		Find(&doctors).Error

	return doctors, err
}

// SearchDoctorsByService выполняет поиск врачей по названию предоставляемой услуги.
func (r *DoctorPostgres) SearchDoctorsByService(ctx context.Context, serviceQuery string) ([]models.Doctor, error) {
	var doctors []models.Doctor
	searchPattern := "%" + strings.ToLower(serviceQuery) + "%"

	err := r.db.WithContext(ctx).Preload("Specialty").
		Joins("JOIN medical_center.services ON medical_center.services.doctor_id = medical_center.doctors.id").
		Where("LOWER(medical_center.services.name) LIKE ?", searchPattern).
		Distinct().
		Find(&doctors).Error

	return doctors, err
}
