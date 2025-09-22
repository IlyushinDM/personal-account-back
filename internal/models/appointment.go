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
	ResultFileURL        sql.NullString `db:"result_file_url" json:"resultFileURL,omitzero"`
	CreatedAt            time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt            time.Time      `db:"updated_at" json:"updatedAt"`

	// Связанные данные для GORM Preload
	Doctor  Doctor  `gorm:"foreignKey:DoctorID"`
	Service Service `gorm:"foreignKey:ServiceID"`
}

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

// Review представляет отзыв пациента о враче
type Review struct {
	ID          uint64         `gorm:"primarykey" db:"id" json:"id"`
	UserID      uint64         `db:"user_id" json:"userID"`
	DoctorID    uint64         `db:"doctor_id" json:"doctorID"`
	Rating      uint16         `db:"rating" json:"rating"`
	Comment     sql.NullString `db:"comment" json:"comment,omitempty"`
	IsModerated bool           `db:"is_moderated" json:"isModerated"`
	CreatedAt   time.Time      `db:"created_at" json:"createdAt"`
}

// Recommendation представляет DTO для текста рекомендации.
type Recommendation struct {
	Text string `json:"text"`
}

// AvailableDatesResponse представляет DTO для ответа со свободными датами.
type AvailableDatesResponse struct {
	SpecialistID   uint64   `json:"specialistId"`
	Month          string   `json:"month"`
	AvailableDates []string `json:"availableDates"`
}

// AvailableSlotsResponse представляет DTO для ответа со свободными слотами.
type AvailableSlotsResponse struct {
	SpecialistID   uint64   `json:"specialistId"`
	Date           string   `json:"date"`
	AvailableSlots []string `json:"availableSlots"`
}
