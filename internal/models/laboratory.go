package models

import (
	"database/sql"
	"time"
)

// LabAnalysis представляет лабораторный анализ
type LabAnalysis struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64         `gorm:"not null;index" json:"user_id"`
	AppointmentID sql.NullInt64  `gorm:"index" json:"appointment_id,omitempty"` // Анализ может быть без записи? Пока - да
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	AssignedDate  time.Time      `gorm:"type:date;not null" json:"assigned_date"`
	StatusID      uint32         `gorm:"not null;index" json:"status_id"`
	ResultFileURL sql.NullString `gorm:"type:varchar(512)" json:"result_file_url,omitempty"`
	ClinicID      sql.NullInt64  `gorm:"index" json:"clinic_id,omitempty"`

	// Связи
	Status AnalysisStatus `gorm:"foreignKey:StatusID" json:"status,omitempty"`
	Clinic Clinic         `gorm:"foreignKey:ClinicID" json:"clinic,omitempty"`
}
