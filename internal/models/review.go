package models

import (
	"database/sql"
	"time"
)

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
