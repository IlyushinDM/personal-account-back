package models

// City представляет город, где живёт (прописан) клиент
type City struct {
	ID   uint32 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// Department представляет отделение клиники
type Department struct {
	ID   uint32 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// Specialty представляет врачебную специальность
type Specialty struct {
	ID           uint32 `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	DepartmentID uint32 `db:"department_id" json:"departmentID"`
}

// AppointmentStatus представляет статус записи на прием
type AppointmentStatus struct {
	ID   uint32 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// AnalysisStatus представляет статус анализа
type AnalysisStatus struct {
	ID   uint32 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// DepartmentWithSpecialties - это DTO (Data Transfer Object) для сервисного слоя,
// представляющий отделение с вложенным списком его специальностей.
type DepartmentWithSpecialties struct {
	ID          uint32      `json:"id"`
	Name        string      `json:"name"`
	Specialties []Specialty `json:"specialties"`
}
