// Package models определяет основные структуры данных (сущности) приложения,
// такие как User, Doctor и Appointment. Эти структуры используются на всех слоях,
// от взаимодействия с базой данных до формирования JSON-ответов.
package models

import (
	"database/sql"
	"time"
)

// User представляет пользователя системы
type User struct {
	ID           uint64         `gorm:"primarykey" db:"id" json:"id"`
	Phone        string         `gorm:"unique" db:"phone" json:"phone"`
	PasswordHash string         `db:"password_hash" json:"-"`
	GosuslugiID  sql.NullString `gorm:"unique" db:"gosuslugi_id" json:"gosuslugiID,omitempty"`
	IsActive     bool           `db:"is_active" json:"isActive"`
	CreatedAt    time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updatedAt"`
}

// UserProfile содержит расширенную информацию о пользователе
type UserProfile struct {
	ID         uint64         `gorm:"primarykey" db:"id" json:"id"`
	UserID     uint64         `gorm:"uniqueIndex" db:"user_id" json:"userID"`
	FirstName  string         `db:"first_name" json:"firstName"`
	LastName   string         `db:"last_name" json:"lastName"`
	Patronymic sql.NullString `db:"patronymic" json:"patronymic,omitempty"`
	BirthDate  time.Time      `db:"birth_date" json:"birthDate"`
	Gender     string         `db:"gender" json:"gender"`
	CityID     uint32         `db:"city_id" json:"cityID"`
	Email      sql.NullString `gorm:"unique" db:"email" json:"email,omitempty"`
	AvatarURL  sql.NullString `db:"avatar_url" json:"avatarURL,omitempty"`
}
