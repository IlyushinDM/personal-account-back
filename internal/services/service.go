// Package services содержит бизнес-логику приложения.
package services

import (
	"context"
	"mime/multipart"
	"time"

	"lk/internal/models"
	"lk/internal/repository"
	"lk/internal/storage"
)

// Authorization определяет методы для регистрации и входа пользователя.
type Authorization interface {
	CreateUser(ctx context.Context, phone, password, fullName,
		gender, birthDateStr string, cityID uint32) (map[string]string, error)
	GenerateToken(ctx context.Context, phone, password string) (map[string]string, error)
	ParseToken(token string) (uint64, error)
	RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error)
	Logout(ctx context.Context, refreshToken string) error
	ForgotPassword(ctx context.Context, phone string) error
	ResetPassword(ctx context.Context, phone, code, newPassword string) error

	// TODO: FR-1.3, FR-1.4 - Добавить методы для работы с Госуслугами
	// GetGosuslugiAuthURL(state string) (string, error)
	// AuthorizeGosuslugi(ctx context.Context, gosuslugiCode, state string) (map[string]string, error)
}

// UserService определяет методы для работы с данными пользователя.
type UserService interface {
	GetFullUserProfile(ctx context.Context, userID uint64) (models.UserProfile, []models.Appointment, error)
	UpdateUserProfile(ctx context.Context, userID uint64, input models.UserProfile) (models.UserProfile, error)
	UpdateAvatar(ctx context.Context, userID uint64, fileHeader *multipart.FileHeader) (string, error)
}

// DoctorService определяет методы для работы с информацией о врачах.
type DoctorService interface {
	GetDoctorByID(ctx context.Context, doctorID uint64) (models.Doctor, error)
	GetDoctorsBySpecialty(ctx context.Context, specialtyID uint32, params models.PaginationParams) (
		models.PaginatedDoctorsResponse, error)
	SearchDoctors(ctx context.Context, query string) ([]models.Doctor, error)
	SearchDoctorsByService(ctx context.Context, serviceQuery string) ([]models.Doctor, error)
	GetSpecialistRecommendations(ctx context.Context, doctorID uint64) (models.Recommendation, error)
}

// DirectoryService для работы со справочниками.
type DirectoryService interface {
	GetAllDepartmentsWithSpecialties(ctx context.Context) ([]models.DepartmentWithSpecialties, error)
	GetSpecialties(ctx context.Context, departmentID *uint32) ([]models.Specialty, error)
}

// AppointmentService определяет методы для работы с записями на прием.
type AppointmentService interface {
	CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error)
	GetUserAppointments(ctx context.Context, userID uint64) ([]models.Appointment, error)
	CancelAppointment(ctx context.Context, userID, appointmentID uint64) error
	GetAvailableDates(ctx context.Context, doctorID, serviceID uint64, month string) (
		models.AvailableDatesResponse, error)
	GetAvailableSlots(ctx context.Context, doctorID, serviceID uint64, date string) (
		models.AvailableSlotsResponse, error)
	GetUpcomingForUser(ctx context.Context, userID uint64) ([]models.Appointment, error)
}

// InfoService определяет методы для работы с общей информацией.
type InfoService interface {
	GetServiceRecommendations(ctx context.Context, serviceID uint64) (models.Recommendation, error)
	GetClinicInfo(ctx context.Context) (models.ClinicInfo, error)
	GetLegalDocuments(ctx context.Context) ([]models.LegalDocument, error)
}

// PrescriptionService определяет методы для работы с назначениями.
type PrescriptionService interface {
	GetActiveForUser(ctx context.Context, userID uint64) ([]models.Prescription, error)
	ArchiveForUser(ctx context.Context, userID, prescriptionID uint64) error
}

// MedicalCardService определяет методы для работы с медкартой.
type MedicalCardService interface {
	GetVisits(ctx context.Context, userID uint64, params models.PaginationParams) (
		models.PaginatedVisitsResponse, error)
	GetAnalyses(ctx context.Context, userID uint64, status *string) ([]models.LabAnalysis, error)
	GetArchivedPrescriptions(ctx context.Context, userID uint64) ([]models.Prescription, error)
	GetSummary(ctx context.Context, userID uint64) (models.MedicalCardSummary, error)
	ArchivePrescription(ctx context.Context, userID, prescriptionID uint64) error
	DownloadFile(ctx context.Context, userID, fileID uint64) ([]byte, string, error)
}

// Service - это контейнер для всех сервисов приложения.
type Service struct {
	Authorization Authorization
	User          UserService
	Doctor        DoctorService
	Appointment   AppointmentService
	Directory     DirectoryService
	Info          InfoService
	Prescription  PrescriptionService
	MedicalCard   MedicalCardService
}

// ServiceDependencies содержит все зависимости, необходимые для создания сервисов.
type ServiceDependencies struct {
	Repos      *repository.Repository
	Storage    storage.FileStorage
	Location   *time.Location
	SigningKey string
	TokenTTL   time.Duration
}

// NewService создает новый экземпляр главного сервиса, инициализируя все реализации.
func NewService(deps ServiceDependencies) *Service {
	authService := NewAuthService(
		deps.Repos.User,
		deps.Repos.Token,
		deps.Repos.Cache,
		deps.Repos.Transactor,
		deps.SigningKey,
		deps.TokenTTL,
	)

	return &Service{
		Authorization: authService,
		User:          NewUserService(deps.Repos.User, deps.Repos.Appointment, deps.Storage),
		Doctor:        NewDoctorService(deps.Repos.Doctor),
		Appointment:   NewAppointmentService(deps.Repos.Appointment, deps.Repos.Doctor, deps.Location),
		Directory:     NewDirectoryService(deps.Repos.Directory),
		Info:          NewInfoService(deps.Repos.Service, deps.Repos.Info),
		Prescription:  NewPrescriptionService(deps.Repos.Prescription),
		MedicalCard:   NewMedicalCardService(deps.Repos.MedicalCard, deps.Repos.Prescription, deps.Storage),
	}
}
