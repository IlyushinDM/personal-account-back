package models

// City представляет город, где живёт (прописан) клиент
type City struct {
	ID   uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(100);not null;unique" json:"name"`
}

// Department представляет отделение клиники
type Department struct {
	ID   uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(100);not null;unique" json:"name"`
}

// Specialty представляет врачебную специальность
type Specialty struct {
	ID           uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"type:varchar(150);not null;unique" json:"name"`
	DepartmentID uint32 `gorm:"not null" json:"department_id"`

	// Связи
	Department Department `gorm:"foreignKey:DepartmentID" json:"department"`
}

// AppointmentStatus представляет статус записи на прием
type AppointmentStatus struct {
	ID   uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(50);not null;unique" json:"name"`
}

// AnalysisStatus представляет статус анализа
type AnalysisStatus struct {
	ID   uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(50);not null;unique" json:"name"`
}
