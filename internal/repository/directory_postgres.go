package repository

import (
	"context"

	"lk/internal/models"

	"github.com/jmoiron/sqlx"
)

// DirectoryPostgres реализует DirectoryRepository для PostgreSQL.
type DirectoryPostgres struct {
	db *sqlx.DB
}

// NewDirectoryPostgres создает новый экземпляр репозитория для справочников.
func NewDirectoryPostgres(db *sqlx.DB) *DirectoryPostgres {
	return &DirectoryPostgres{db: db}
}

// GetAllDepartments возвращает список всех отделений клиники.
func (r *DirectoryPostgres) GetAllDepartments(ctx context.Context) ([]models.Department, error) {
	var departments []models.Department
	query := "SELECT * FROM medical_center.departments ORDER BY name"
	err := r.db.SelectContext(ctx, &departments, query)
	return departments, err
}

// GetAllSpecialties возвращает список всех врачебных специальностей.
// * Если departmentID не является nil, фильтрует по ID отделения.
func (r *DirectoryPostgres) GetAllSpecialties(ctx context.Context, departmentID *uint32) ([]models.Specialty, error) {
	var specialties []models.Specialty
	query := "SELECT * FROM medical_center.specialties"
	args := []interface{}{}

	if departmentID != nil {
		query += " WHERE department_id=$1"
		args = append(args, *departmentID)
	}

	query += " ORDER BY name"

	err := r.db.SelectContext(ctx, &specialties, query, args...)
	return specialties, err
}
