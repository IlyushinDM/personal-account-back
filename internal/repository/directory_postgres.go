package repository

import (
	"context"

	"lk/internal/models"

	"gorm.io/gorm"
)

// DirectoryPostgres реализует DirectoryRepository для PostgreSQL.
type DirectoryPostgres struct {
	db *gorm.DB
}

// NewDirectoryPostgres создает новый экземпляр репозитория для справочников.
func NewDirectoryPostgres(db *gorm.DB) *DirectoryPostgres {
	return &DirectoryPostgres{db: db}
}

// GetAllDepartments возвращает список всех отделений клиники.
func (r *DirectoryPostgres) GetAllDepartments(ctx context.Context) ([]models.Department, error) {
	var departments []models.Department
	err := r.db.WithContext(ctx).Order("name").Find(&departments).Error
	return departments, err
}

// GetAllSpecialties возвращает список всех врачебных специальностей.
// * Если departmentID не является nil, фильтрует по ID отделения.
func (r *DirectoryPostgres) GetAllSpecialties(ctx context.Context, departmentID *uint32) ([]models.Specialty, error) {
	var specialties []models.Specialty
	query := r.db.WithContext(ctx).Order("name")

	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	err := query.Find(&specialties).Error
	return specialties, err
}
