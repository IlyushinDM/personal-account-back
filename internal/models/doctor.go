package models

import (
	"database/sql"
	"time"
)

// Doctor представляет профиль врача
type Doctor struct {
	ID              uint64         `gorm:"primarykey" db:"id" json:"id"`
	FirstName       string         `db:"first_name" json:"firstName"`
	LastName        string         `db:"last_name" json:"lastName"`
	Patronymic      sql.NullString `db:"patronymic" json:"patronymic,omitempty"`
	SpecialtyID     uint32         `db:"specialty_id" json:"specialtyID"`
	ExperienceYears uint16         `db:"experience_years" json:"experienceYears"`
	Rating          float32        `db:"rating" json:"rating"`
	ReviewCount     uint32         `db:"review_count" json:"reviewCount"`
	AvatarURL       sql.NullString `db:"avatar_url" json:"avatarURL,omitempty"`
	Recommendations sql.NullString `json:"recommendations,omitempty"`
	CreatedAt       time.Time      `db:"created_at" json:"createdAt"`
	Specialty       Specialty      `gorm:"foreignKey:SpecialtyID" db:"specialty" json:"specialty"`
}

// DoctorEducation описывает образование врача
type DoctorEducation struct {
	ID          uint64 `gorm:"primarykey" db:"id" json:"id"`
	DoctorID    uint64 `db:"doctor_id" json:"doctorID"`
	Institution string `db:"institution" json:"institution"`
	Specialty   string `db:"specialty" json:"specialty"`
	StartYear   uint16 `db:"start_year" json:"startYear"`
	EndYear     uint16 `db:"end_year" json:"endYear"`
}

// DoctorResidency описывает ординатуру врача
type DoctorResidency struct {
	ID          uint64 `gorm:"primarykey" db:"id" json:"id"`
	DoctorID    uint64 `db:"doctor_id" json:"doctorID"`
	Institution string `db:"institution" json:"institution"`
	Specialty   string `db:"specialty" json:"specialty"`
	StartYear   uint16 `db:"start_year" json:"startYear"`
	EndYear     uint16 `db:"end_year" json:"endYear"`
}

// DoctorCourse описывает курсы повышения квалификации
type DoctorCourse struct {
	ID         uint64 `gorm:"primarykey" db:"id" json:"id"`
	DoctorID   uint64 `db:"doctor_id" json:"doctorID"`
	CourseName string `db:"course_name" json:"courseName"`
	Year       uint16 `db:"year" json:"year"`
}

// DoctorCertificate описывает сертификаты врача
type DoctorCertificate struct {
	ID         uint64         `gorm:"primarykey" db:"id" json:"id"`
	DoctorID   uint64         `db:"doctor_id" json:"doctorID"`
	CertName   string         `db:"cert_name" json:"certName"`
	CertNumber sql.NullString `db:"cert_number" json:"certNumber,omitempty"`
}

// DoctorSkill описывает профессиональные навыки врача
type DoctorSkill struct {
	ID        uint64 `gorm:"primarykey" db:"id" json:"id"`
	DoctorID  uint64 `db:"doctor_id" json:"doctorID"`
	SkillName string `db:"skill_name" json:"skillName"`
}

// DoctorSpecialization описывает узкие специализации врача
type DoctorSpecialization struct {
	ID       uint64 `gorm:"primarykey" db:"id" json:"id"`
	DoctorID uint64 `db:"doctor_id" json:"doctorID"`
	Area     string `db:"area" json:"area"`
}

// PaginatedDoctorsResponse - это DTO для возврата пагинированного списка врачей.
type PaginatedDoctorsResponse struct {
	Total int64    `json:"total" example:"25"`
	Items []Doctor `json:"items"`
}
