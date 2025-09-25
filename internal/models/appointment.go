package models

import (
	"database/sql"
	"time"
)

// Appointment представляет запись на прием к врачу
type Appointment struct {
	ID                   uint64         `gorm:"primarykey" db:"id" json:"id"`
	UserID               uint64         `db:"user_id" json:"userID"`
	DoctorID             uint64         `db:"doctor_id" json:"doctorID"`
	ServiceID            uint64         `db:"service_id" json:"serviceID"`
	ClinicID             uint64         `db:"clinic_id" json:"clinicID"`
	AppointmentDate      time.Time      `db:"appointment_date" json:"appointmentDate"`
	AppointmentTime      string         `db:"appointment_time" json:"appointmentTime"`
	StatusID             uint32         `db:"status_id" json:"statusID"`
	PriceAtBooking       float64        `db:"price_at_booking" json:"priceAtBooking"`
	IsDMS                bool           `db:"is_dms" json:"isDMS"`
	PreVisitInstructions sql.NullString `db:"pre_visit_instructions" json:"preVisitInstructions,omitzero"`
	Diagnosis            sql.NullString `db:"diagnosis" json:"diagnosis,omitzero"`
	Recommendations      sql.NullString `db:"recommendations" json:"recommendations,omitzero"`
	ResultFileURL        sql.NullString `db:"result_file_url" json:"-"`
	CreatedAt            time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt            time.Time      `db:"updated_at" json:"updatedAt"`
	Doctor               Doctor         `gorm:"foreignKey:DoctorID"`
	Service              Service        `gorm:"foreignKey:ServiceID"`
}

func (Appointment) TableName() string {
	return "medical_center.appointments"
}
