package models

import (
	"database/sql"
	"time"
)

// Doctor представляет профиль врача
type Doctor struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName       string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName        string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Patronymic      sql.NullString `gorm:"type:varchar(100)" json:"patronymic,omitzero"`
	SpecialtyID     uint32         `gorm:"not null" json:"specialty_id"`
	ExperienceYears uint16         `gorm:"not null" json:"experience_years"`
	Rating          float32        `gorm:"type:numeric(3,2);not null;default:0.00" json:"rating"`
	ReviewCount     uint32         `gorm:"not null;default:0" json:"review_count"`
	AvatarURL       sql.NullString `gorm:"type:varchar(512)" json:"avatar_url,omitzero"`
	CreatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Связи для полной загрузки профиля
	Specialty       Specialty              `gorm:"foreignKey:SpecialtyID" json:"specialty"`
	Clinics         []Clinic               `gorm:"many2many:doctor_clinics" json:"clinics,omitempty"`
	Services        []Service              `gorm:"foreignKey:DoctorID" json:"services,omitempty"`
	Education       []DoctorEducation      `gorm:"foreignKey:DoctorID" json:"education,omitempty"`
	Residency       []DoctorResidency      `gorm:"foreignKey:DoctorID" json:"residency,omitempty"`
	Courses         []DoctorCourse         `gorm:"foreignKey:DoctorID" json:"courses,omitempty"`
	Certificates    []DoctorCertificate    `gorm:"foreignKey:DoctorID" json:"certificates,omitempty"`
	Skills          []DoctorSkill          `gorm:"foreignKey:DoctorID" json:"skills,omitempty"`
	Specializations []DoctorSpecialization `gorm:"foreignKey:DoctorID" json:"specializations,omitempty"`
	Reviews         []Review               `gorm:"foreignKey:DoctorID" json:"reviews,omitempty"`
}

// DoctorEducation описывает образование врача
type DoctorEducation struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	DoctorID    uint64 `gorm:"not null" json:"doctor_id"`
	Institution string `gorm:"type:varchar(255);not null" json:"institution"`
	Specialty   string `gorm:"type:varchar(255);not null" json:"specialty"`
	StartYear   uint16 `gorm:"not null" json:"start_year"`
	EndYear     uint16 `gorm:"not null" json:"end_year"`
}

// DoctorResidency описывает ординатуру врача
type DoctorResidency struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	DoctorID    uint64 `gorm:"not null" json:"doctor_id"`
	Institution string `gorm:"type:varchar(255);not null" json:"institution"`
	Specialty   string `gorm:"type:varchar(255);not null" json:"specialty"`
	StartYear   uint16 `gorm:"not null" json:"start_year"`
	EndYear     uint16 `gorm:"not null" json:"end_year"`
}

// DoctorCourse описывает курсы повышения квалификации
type DoctorCourse struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	DoctorID   uint64 `gorm:"not null" json:"doctor_id"`
	CourseName string `gorm:"type:varchar(255);not null" json:"course_name"`
	Year       uint16 `gorm:"not null" json:"year"`
}

// DoctorCertificate описывает сертификаты врача
type DoctorCertificate struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	DoctorID   uint64         `gorm:"not null" json:"doctor_id"`
	CertName   string         `gorm:"type:varchar(255);not null" json:"cert_name"`
	CertNumber sql.NullString `gorm:"type:varchar(100)" json:"cert_number,omitzero"`
}

// DoctorSkill описывает профессиональные навыки врача
type DoctorSkill struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	DoctorID  uint64 `gorm:"not null" json:"doctor_id"`
	SkillName string `gorm:"type:varchar(255);not null" json:"skill_name"`
}

// DoctorSpecialization описывает узкие специализации врача
type DoctorSpecialization struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	DoctorID uint64 `gorm:"not null" json:"doctor_id"`
	Area     string `gorm:"type:varchar(255);not null" json:"area"`
}
