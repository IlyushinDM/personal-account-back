package models

import (
	"database/sql"
	"time"
)

// LabAnalysis представляет лабораторный анализ
type LabAnalysis struct {
	ID             uint64         `gorm:"primarykey" db:"id" json:"id"`
	UserID         uint64         `db:"user_id" json:"userID"`
	AppointmentID  sql.NullInt64  `db:"appointment_id" json:"appointmentID,omitempty"`
	Name           string         `db:"name" json:"name"`
	AssignedDate   time.Time      `db:"assigned_date" json:"assignedDate"`
	StatusID       uint32         `db:"status_id" json:"statusID"`
	ResultFileURL  sql.NullString `db:"result_file_url" json:"-"`
	ResultFileName sql.NullString `db:"result_file_name" json:"resultFileName,omitempty"`
	ClinicID       sql.NullInt64  `db:"clinic_id" json:"clinicID,omitempty"`
}

func (LabAnalysis) TableName() string {
	return "medical_center.labanalyses"
}
