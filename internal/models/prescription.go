package models

import (
	"database/sql"
	"time"
)

// Prescription представляет назначение/рецепт от врача
type Prescription struct {
	ID            uint64       `gorm:"primarykey" db:"id" json:"id"`
	AppointmentID uint64       `db:"appointment_id" json:"appointmentID"`
	UserID        uint64       `db:"user_id" json:"userID"`
	DoctorID      uint64       `db:"doctor_id" json:"doctorID"`
	Content       string       `db:"content" json:"content"`
	Status        string       `db:"status" json:"status"`
	CompletedAt   sql.NullTime `db:"completed_at" json:"completedAt,omitempty"`
	CreatedAt     time.Time    `db:"created_at" json:"createdAt"`
	ArchivedDate  sql.NullTime `json:"archivedDate,omitempty"`
}
