// Package repository определяет слой доступа к данным (DAL) приложения.
package repository

import (
	"context"
	"time"

	"lk/internal/models"

	"github.com/redis/go-redis/v9"
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
	UpdatePassword(ctx context.Context, userID uint64, newPasswordHash string) error
}

// TokenRepository определяет методы для работы с refresh-токенами.
type TokenRepository interface {
	Create(ctx context.Context, token models.RefreshToken) error
	GetByUserID(ctx context.Context, userID uint64) (models.RefreshToken, error)
	Delete(ctx context.Context, userID uint64) error
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
// ИНТЕРФЕЙС ОБНОВЛЕН
type AppointmentRepository interface {
	CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error)
	GetAppointmentsByUserID(ctx context.Context, userID uint64) ([]models.Appointment, error)
	GetUpcomingAppointmentsByUserID(ctx context.Context, userID uint64) ([]models.Appointment, error)
	GetAppointmentByID(ctx context.Context, appointmentID uint64) (models.Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, appointmentID uint64, statusID uint32) error

	// Новые методы для работы с реальным расписанием
	GetAvailableDatesForMonth(ctx context.Context, doctorID uint64, month time.Time) ([]time.Time, error)
	GetDoctorScheduleForDate(ctx context.Context, doctorID uint64, date time.Time) (models.Schedule, error)
	GetServiceDurationMinutes(ctx context.Context, serviceID uint64) (uint16, error)
	GetAppointmentsByDoctorAndDate(ctx context.Context, doctorID uint64, date time.Time) ([]models.Appointment, error)
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

// MedicalCardRepository определяет методы для работы с данными медкарты.
type MedicalCardRepository interface {
	GetCompletedVisits(ctx context.Context, userID uint64, params models.PaginationParams) (
		[]models.Appointment, int64, error)
	GetAnalysesByUserID(ctx context.Context, userID uint64, status *string) ([]models.LabAnalysis, error)
	GetArchivedPrescriptionsByUserID(ctx context.Context, userID uint64) ([]models.Prescription, error)
	GetSummaryInfo(ctx context.Context, userID uint64) (models.MedicalCardSummary, error)
	ArchivePrescription(ctx context.Context, userID, prescriptionID uint64) error
	GetAnalysisByID(ctx context.Context, analysisID uint64) (models.LabAnalysis, error)
}

// CacheRepository определяет интерфейс для работы с key-value хранилищем (кэшем).
type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

// Repository - это контейнер для всех репозиториев приложения.
type Repository struct {
	User         UserRepository
	Token        TokenRepository
	Doctor       DoctorRepository
	Appointment  AppointmentRepository
	Directory    DirectoryRepository
	Info         InfoRepository
	Prescription PrescriptionRepository
	Service      ServiceRepository
	MedicalCard  MedicalCardRepository
	Cache        CacheRepository
	Transactor
}

// NewRepository создает новый экземпляр главного репозитория.
func NewRepository(db *gorm.DB, redisClient *redis.Client) *Repository {
	return &Repository{
		User:         NewUserPostgres(db),
		Token:        NewTokenPostgres(db),
		Doctor:       NewDoctorPostgres(db),
		Appointment:  NewAppointmentPostgres(db),
		Directory:    NewDirectoryPostgres(db),
		Info:         NewInfoPostgres(db),
		Prescription: NewPrescriptionPostgres(db),
		Service:      NewServicePostgres(db),
		MedicalCard:  NewMedicalCardPostgres(db),
		Cache:        NewCacheRedis(redisClient),
		Transactor:   NewTransactor(db),
	}
}
