// Package repository определяет слой доступа к данным (DAL) приложения.
package repository

import (
	"context"
	"time"

	"lk/internal/models"

	"gorm.io/gorm"
)

// Transactor определяет интерфейс для управления транзакциями.
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

// UserRepository определяет методы для работы с пользователями.
type UserRepository interface {
	CreateUser(ctx context.Context, tx *gorm.DB, user models.User) (uint64, error)
	GetUserByPhone(ctx context.Context, phone string) (models.User, error)
	GetUserByID(ctx context.Context, id uint64) (models.User, error)
	CreateUserProfile(ctx context.Context, tx *gorm.DB, profile models.UserProfile) (uint64, error)
	GetUserProfileByUserID(ctx context.Context, userID uint64) (models.UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile models.UserProfile) (models.UserProfile, error)
	UpdateAvatar(ctx context.Context, userID uint64, avatarURL string) error
}

// DoctorRepository определяет методы для работы с врачами.
type DoctorRepository interface {
	GetDoctorByID(ctx context.Context, id uint64) (models.Doctor, error)
	GetDoctorsBySpecialty(ctx context.Context, specialtyID uint32, params models.PaginationParams) (
		[]models.Doctor, int64, error)
	SearchDoctors(ctx context.Context, query string) ([]models.Doctor, error)
	SearchDoctorsByService(ctx context.Context, serviceQuery string) ([]models.Doctor, error)
	GetSpecialistRecommendations(ctx context.Context, doctorID uint64) (string, error)
}

// AppointmentRepository определяет методы для работы с записями на прием.
type AppointmentRepository interface {
	CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error)
	GetAppointmentsByUserID(ctx context.Context, userID uint64) ([]models.Appointment, error)
	GetUpcomingAppointmentsByUserID(ctx context.Context, userID uint64) ([]models.Appointment, error)
	GetAppointmentByID(ctx context.Context, appointmentID uint64) (models.Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, appointmentID uint64, statusID uint32) error
	GetAvailableDates(ctx context.Context, doctorID uint64, month time.Time) ([]time.Time, error)
	GetAvailableSlots(ctx context.Context, doctorID uint64, date time.Time) ([]string, error)
}

// DirectoryRepository определяет методы для работы со справочниками.
type DirectoryRepository interface {
	GetAllDepartments(ctx context.Context) ([]models.Department, error)
	GetAllSpecialties(ctx context.Context, departmentID *uint32) ([]models.Specialty, error)
}

// InfoRepository определяет методы для работы с общей информацией.
type InfoRepository interface {
	GetClinicInfo(ctx context.Context) (models.ClinicInfo, error)
	GetLegalDocuments(ctx context.Context) ([]models.LegalDocument, error)
}

// PrescriptionRepository определяет методы для работы с назначениями.
type PrescriptionRepository interface {
	GetActiveByUserID(ctx context.Context, userID uint64) ([]models.Prescription, error)
	Archive(ctx context.Context, userID, prescriptionID uint64) error
	GetByID(ctx context.Context, prescriptionID uint64) (models.Prescription, error)
}

// ServiceRepository определяет методы для работы с услугами.
type ServiceRepository interface {
	GetServiceRecommendations(ctx context.Context, serviceID uint64) (string, error)
}

// Repository - это контейнер для всех репозиториев приложения.
type Repository struct {
	User         UserRepository
	Doctor       DoctorRepository
	Appointment  AppointmentRepository
	Directory    DirectoryRepository
	Info         InfoRepository
	Prescription PrescriptionRepository
	Service      ServiceRepository
	MedicalCard  MedicalCardRepository
	Transactor
}

// NewRepository создает новый экземпляр главного репозитория.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:         NewUserPostgres(db),
		Doctor:       NewDoctorPostgres(db),
		Appointment:  NewAppointmentPostgres(db),
		Directory:    NewDirectoryPostgres(db),
		Info:         NewInfoPostgres(db),
		Prescription: NewPrescriptionPostgres(db),
		Service:      NewServicePostgres(db),
		MedicalCard:  NewMedicalCardPostgres(db),
		Transactor:   NewTransactor(db),
	}
}

// MedicalCardRepository определяет методы для работы с данными медкарты.
type MedicalCardRepository interface {
	// FR-4.1
	GetCompletedVisits(ctx context.Context, userID uint64, params models.PaginationParams) (
		[]models.Appointment, int64, error)
	// FR-4.2
	GetAnalysesByUserID(ctx context.Context, userID uint64, status *string) ([]models.LabAnalysis, error)
	// FR-4.3
	GetArchivedPrescriptionsByUserID(ctx context.Context, userID uint64) ([]models.Prescription, error)
	// FR-4.6
	GetSummaryInfo(ctx context.Context, userID uint64) (models.MedicalCardSummary, error)
	// FR-4.7
	ArchivePrescription(ctx context.Context, userID, prescriptionID uint64) error
}
