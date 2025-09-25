package models

// City представляет город, где живёт (прописан) клиент
type City struct {
	ID   uint32 `gorm:"primarykey" db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (City) TableName() string {
	return "medical_center.cities"
}

// Department представляет отделение клиники
type Department struct {
	ID   uint32 `gorm:"primarykey" db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (Department) TableName() string {
	return "medical_center.departments"
}

// Specialty представляет врачебную специальность
type Specialty struct {
	ID           uint32 `gorm:"primarykey" db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	DepartmentID uint32 `db:"department_id" json:"departmentID"`
}

func (Specialty) TableName() string {
	return "medical_center.specialties"
}

// AppointmentStatus представляет статус записи на прием
type AppointmentStatus struct {
	ID   uint32 `gorm:"primarykey" db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (AppointmentStatus) TableName() string {
	return "medical_center.appointmentstatuses"
}

// AnalysisStatus представляет статус анализа
type AnalysisStatus struct {
	ID   uint32 `gorm:"primarykey" db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (AnalysisStatus) TableName() string {
	return "medical_center.analysisstatuses"
}

// DepartmentWithSpecialties - это DTO (Data Transfer Object) для сервисного слоя,
// представляющий отделение с вложенным списком его специальностей.
type DepartmentWithSpecialties struct {
	ID          uint32      `json:"id"`
	Name        string      `json:"name"`
	Specialties []Specialty `json:"specialties"`
}
