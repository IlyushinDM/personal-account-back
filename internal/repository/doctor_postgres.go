package repository

import (
	"context"
	"fmt"
	"strings"

	"lk/internal/models"

	"github.com/jmoiron/sqlx"
)

// DoctorPostgres реализует интерфейс DoctorRepository для PostgreSQL.
type DoctorPostgres struct {
	db *sqlx.DB
}

// NewDoctorPostgres создает новый экземпляр репозитория для врачей.
func NewDoctorPostgres(db *sqlx.DB) *DoctorPostgres {
	return &DoctorPostgres{db: db}
}

// GetDoctorByID получает информацию о враче по его ID.
func (r *DoctorPostgres) GetDoctorByID(ctx context.Context, id uint64) (models.Doctor, error) {
	var doctor models.Doctor
	query := `
		SELECT d.*, s.name as "specialty.name"
		FROM medical_center.doctors d
		LEFT JOIN medical_center.specialties s ON d.specialty_id = s.id
		WHERE d.id=$1`

	err := r.db.GetContext(ctx, &doctor, query, id)
	return doctor, err
}

// GetSpecialistRecommendations возвращает фиктивный текст рекомендаций для специалиста.
// ! Это mock-реализация. В реальном приложении здесь будет запрос к полю в таблице doctors.
func (r *DoctorPostgres) GetSpecialistRecommendations(ctx context.Context, doctorID uint64) (string, error) {
	// Просто проверяем, что доктор существует
	_, err := r.GetDoctorByID(ctx, doctorID)
	if err != nil {
		return "", err
	}
	return "Пожалуйста, приходите за 10 минут до начала приема и возьмите с собой все предыдущие медицинские заключения.", nil
}

// GetDoctorsBySpecialty получает список врачей по ID их специальности с пагинацией и сортировкой.
// Возвращает список врачей и общее количество врачей по данной специальности.
func (r *DoctorPostgres) GetDoctorsBySpecialty(ctx context.Context, specialtyID uint32, params models.PaginationParams) ([]models.Doctor, int, error) {
	var doctors []models.Doctor
	var total int

	// Белый список для колонок сортировки для предотвращения SQL-инъекций
	allowedSortBy := map[string]string{
		"rating":           "d.rating",
		"experience":       "d.experience_years",
		"experience_years": "d.experience_years",
		"name":             "d.last_name",
	}
	orderByColumn, ok := allowedSortBy[params.SortBy]
	if !ok {
		orderByColumn = "d.rating" // Сортировка по умолчанию
	}

	sortOrder := "DESC"
	if strings.ToUpper(params.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	baseQuery := `
		FROM medical_center.doctors d
		LEFT JOIN medical_center.specialties s ON d.specialty_id = s.id
		WHERE d.specialty_id=$1`

	// Получаем общее количество
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.GetContext(ctx, &total, countQuery, specialtyID); err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []models.Doctor{}, 0, nil
	}

	// Получаем пагинированные данные
	dataQuery := fmt.Sprintf(`
		SELECT d.*, s.name as "specialty.name"
		%s
		ORDER BY %s %s
		LIMIT $2 OFFSET $3`, baseQuery, orderByColumn, sortOrder)

	offset := (params.Page - 1) * params.Limit
	err := r.db.SelectContext(ctx, &doctors, dataQuery, specialtyID, params.Limit, offset)
	return doctors, total, err
}

// SearchDoctors выполняет поиск врачей по ФИО или названию специальности.
func (r *DoctorPostgres) SearchDoctors(ctx context.Context, searchQuery string) ([]models.Doctor, error) {
	var doctors []models.Doctor
	searchPattern := "%" + searchQuery + "%"

	query := `
		SELECT d.*, s.name as "specialty.name"
		FROM medical_center.doctors d
		LEFT JOIN medical_center.specialties s ON d.specialty_id = s.id
		WHERE CONCAT_WS(' ', d.last_name, d.first_name, d.patronymic) ILIKE $1 OR s.name ILIKE $1`

	err := r.db.SelectContext(ctx, &doctors, query, searchPattern)
	return doctors, err
}

// SearchDoctorsByService выполняет поиск врачей по названию предоставляемой услуги.
func (r *DoctorPostgres) SearchDoctorsByService(ctx context.Context, serviceQuery string) ([]models.Doctor, error) {
	var doctors []models.Doctor
	searchPattern := "%" + serviceQuery + "%"

	query := `
		SELECT DISTINCT d.*, s.name as "specialty.name"
		FROM medical_center.doctors d
		LEFT JOIN medical_center.specialties s ON d.specialty_id = s.id
		JOIN medical_center.services svc ON d.id = svc.doctor_id
		WHERE svc.name ILIKE $1`

	err := r.db.SelectContext(ctx, &doctors, query, searchPattern)
	return doctors, err
}
