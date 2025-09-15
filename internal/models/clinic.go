package models

import "database/sql"

// Clinic представляет медицинский центр
type Clinic struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"type:varchar(150);not null" json:"name"`
	Address   string `gorm:"type:varchar(512);not null" json:"address"`
	WorkHours string `gorm:"type:varchar(100);not null" json:"work_hours"`
	Phone     string `gorm:"type:varchar(20);not null" json:"phone"`

	// Связи
	Doctors []Doctor `gorm:"many2many:doctor_clinics" json:"doctors,omitempty"`
}

// Service представляет медицинскую услугу, оказываемую врачом
type Service struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string         `gorm:"type:varchar(255);not null" json:"name"`
	Price           float64        `gorm:"type:numeric(10,2);not null" json:"price"`
	DurationMinutes uint16         `gorm:"not null" json:"duration_minutes"`
	Description     sql.NullString `gorm:"type:text" json:"description,omitzero"`
	DoctorID        uint64         `gorm:"not null;index" json:"doctor_id"`
}

// DoctorClinic связывает врачей и клиники
type DoctorClinic struct {
	DoctorID uint64 `gorm:"primaryKey" json:"doctor_id"`
	ClinicID uint64 `gorm:"primaryKey" json:"clinic_id"`
}
