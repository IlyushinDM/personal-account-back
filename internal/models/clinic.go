package models

import "database/sql"

// Clinic представляет медицинский центр
type Clinic struct {
	ID        uint64 `gorm:"primarykey" db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Address   string `db:"address" json:"address"`
	WorkHours string `db:"work_hours" json:"workHours"`
	Phone     string `db:"phone" json:"phone"`
}

// Service представляет медицинскую услугу, оказываемую врачом
type Service struct {
	ID              uint64         `gorm:"primarykey" db:"id" json:"id"`
	Name            string         `db:"name" json:"name"`
	Price           float64        `db:"price" json:"price"`
	DurationMinutes uint16         `db:"duration_minutes" json:"durationMinutes"`
	Description     sql.NullString `db:"description" json:"description,omitempty"`
	Recommendations sql.NullString `json:"recommendations,omitempty"`
	DoctorID        uint64         `db:"doctor_id" json:"doctorID"`
}

// DoctorClinic связывает врачей и клиники
type DoctorClinic struct {
	DoctorID uint64 `gorm:"primaryKey" db:"doctor_id" json:"doctorID"`
	ClinicID uint64 `gorm:"primaryKey" db:"clinic_id" json:"clinicID"`
}
