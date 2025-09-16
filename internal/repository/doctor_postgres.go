package repository

import (
	"context"

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

// GetDoctorsBySpecialty получает список врачей по ID их специальности.
func (r *DoctorPostgres) GetDoctorsBySpecialty(ctx context.Context, specialtyID uint32) ([]models.Doctor, error) {
	var doctors []models.Doctor
	query := `
		SELECT d.*, s.name as "specialty.name"
		FROM medical_center.doctors d
		LEFT JOIN medical_center.specialties s ON d.specialty_id = s.id
		WHERE d.specialty_id=$1`

	err := r.db.SelectContext(ctx, &doctors, query, specialtyID)
	return doctors, err
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
