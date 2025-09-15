package models

import (
	"database/sql"
	"time"
)

// User представляет пользователя системы
type User struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone        string         `gorm:"type:varchar(20);not null;unique" json:"phone"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	GosuslugiID  sql.NullString `gorm:"type:varchar(255);unique" json:"gosuslugi_id,omitzero"`
	IsActive     bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Связи структуры
	Profile      UserProfile   `gorm:"foreignKey:UserID" json:"profile"`
	Appointments []Appointment `gorm:"foreignKey:UserID" json:"appointments,omitempty"`
	Reviews      []Review      `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
	LabAnalyses  []LabAnalysis `gorm:"foreignKey:UserID" json:"lab_analyses,omitempty"`
}

// UserProfile содержит расширенную информацию о пользователе
type UserProfile struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64         `gorm:"not null;unique" json:"user_id"`
	FirstName  string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName   string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Patronymic sql.NullString `gorm:"type:varchar(100)" json:"patronymic,omitzero"`
	BirthDate  time.Time      `gorm:"type:date;not null" json:"birth_date"`
	Gender     string         `gorm:"type:varchar(10);not null" json:"gender"`
	CityID     uint32         `gorm:"not null" json:"city_id"`
	Email      sql.NullString `gorm:"type:varchar(255);unique" json:"email,omitzero"`
	AvatarURL  sql.NullString `gorm:"type:varchar(512)" json:"avatar_url,omitzero"`

	// Связи структуры
	City City `gorm:"foreignKey:CityID" json:"city"`
}
