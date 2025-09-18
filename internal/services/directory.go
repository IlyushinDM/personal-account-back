package services

import (
	"context"

	"lk/internal/models"
	"lk/internal/repository"
)

// directoryService реализует интерфейс DirectoryService.
type directoryService struct {
	repo repository.DirectoryRepository
}

// NewDirectoryService создает новый сервис для работы со справочниками.
func NewDirectoryService(repo repository.DirectoryRepository) DirectoryService {
	return &directoryService{repo: repo}
}

// GetAllDepartmentsWithSpecialties получает иерархию "отделение -> специальности".
func (s *directoryService) GetAllDepartmentsWithSpecialties(ctx context.Context) (
	[]models.DepartmentWithSpecialties, error,
) {
	departments, err := s.repo.GetAllDepartments(ctx)
	if err != nil {
		return nil, err
	}

	// Вызываем метод с nil, чтобы получить ВСЕ специальности без фильтрации.
	specialties, err := s.repo.GetAllSpecialties(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Создаем мапу для быстрого доступа к специальностям по ID отделения.
	specialtiesByDeptID := make(map[uint32][]models.Specialty)
	for _, spec := range specialties {
		specialtiesByDeptID[spec.DepartmentID] = append(specialtiesByDeptID[spec.DepartmentID], spec)
	}

	// Собираем итоговую структуру.
	result := make([]models.DepartmentWithSpecialties, 0, len(departments))
	for _, dept := range departments {
		// Убедимся, что даже если у отделения нет специальностей, поле будет пустым, а не nil
		specs, ok := specialtiesByDeptID[dept.ID]
		if !ok {
			specs = []models.Specialty{}
		}
		result = append(result, models.DepartmentWithSpecialties{
			ID:          dept.ID,
			Name:        dept.Name,
			Specialties: specs,
		})
	}

	return result, nil
}

// GetSpecialties получает список специальностей, опционально фильтруя по отделению.
func (s *directoryService) GetSpecialties(ctx context.Context, departmentID *uint32) ([]models.Specialty, error) {
	return s.repo.GetAllSpecialties(ctx, departmentID)
}
