package models

import "time"

// RefreshToken хранит хеш refresh-токена в базе данных.
type RefreshToken struct {
	ID        uint64    `gorm:"primarykey"`
	UserID    uint64    `gorm:"not null;uniqueIndex"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func (RefreshToken) TableName() string {
	return "medical_center.refresh_tokens"
}
