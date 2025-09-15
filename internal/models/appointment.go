package models

import (
	"database/sql"
	"time"
)

// Appointment представляет запись на прием к врачу
type Appointment struct {
	ID                   uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID               uint64         `gorm:"not null;index" json:"user_id"`
	DoctorID             uint64         `gorm:"not null;index" json:"doctor_id"`
	ServiceID            uint64         `gorm:"not null" json:"service_id"`
	ClinicID             uint64         `gorm:"not null" json:"clinic_id"`
	AppointmentDate      time.Time      `gorm:"type:date;not null" json:"appointment_date"`
	AppointmentTime      string         `gorm:"type:time;not null" json:"appointment_time"`
	StatusID             uint32         `gorm:"not null" json:"status_id"`
	PriceAtBooking       float64        `gorm:"type:numeric(10,2);not null" json:"price_at_booking"`
	IsDMS                bool           `gorm:"not null;default:false" json:"is_dms"`
	PreVisitInstructions sql.NullString `gorm:"type:text" json:"pre_visit_instructions,omitzero"`
	Diagnosis            sql.NullString `gorm:"type:text" json:"diagnosis,omitzero"`
	Recommendations      sql.NullString `gorm:"type:text" json:"recommendations,omitzero"`
	ResultFileURL        sql.NullString `gorm:"type:varchar(512)" json:"result_file_url,omitzero"`
	CreatedAt            time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Связи, необходимые для составления записи
	User    User              `gorm:"foreignKey:UserID" json:"user,omitzero"`
	Doctor  Doctor            `gorm:"foreignKey:DoctorID" json:"doctor,omitzero"`
	Service Service           `gorm:"foreignKey:ServiceID" json:"service,omitzero"`
	Clinic  Clinic            `gorm:"foreignKey:ClinicID" json:"clinic,omitzero"`
	Status  AppointmentStatus `gorm:"foreignKey:StatusID" json:"status,omitzero"`
}

// Prescription представляет назначение/рецепт от врача
type Prescription struct {
	ID            uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	AppointmentID uint64       `gorm:"not null;index" json:"appointment_id"`
	DoctorID      uint64       `gorm:"not null;index" json:"doctor_id"`
	Content       string       `gorm:"type:text;not null" json:"content"`
	Status        string       `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CompletedAt   sql.NullTime `gorm:"null" json:"completed_at,omitzero"`
	CreatedAt     time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Связи
	Appointment Appointment `gorm:"foreignKey:AppointmentID" json:"appointment,omitzero"`
	Doctor      Doctor      `gorm:"foreignKey:DoctorID" json:"doctor,omitzero"`
}

// Review представляет отзыв пациента о враче
type Review struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint64         `gorm:"not null;index" json:"user_id"`
	DoctorID    uint64         `gorm:"not null;index" json:"doctor_id"`
	Rating      uint16         `gorm:"not null" json:"rating"`
	Comment     sql.NullString `gorm:"type:text" json:"comment,omitzero"`
	IsModerated bool           `gorm:"not null;default:false" json:"is_moderated"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Связи
	User User `gorm:"foreignKey:UserID" json:"user,omitzero"`
}
