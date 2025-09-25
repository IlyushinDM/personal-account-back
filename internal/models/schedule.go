package models

import "time"

// Schedule определяет график работы врача в конкретный день.
type Schedule struct {
	ID        uint64    `gorm:"primarykey"`
	DoctorID  uint64    `gorm:"index"`
	Date      time.Time `gorm:"type:date"`
	StartTime time.Time `gorm:"type:time"`
	EndTime   time.Time `gorm:"type:time"`
}

func (Schedule) TableName() string {
	return "medical_center.schedules"
}
